package meme

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"

	"github.com/corona10/goimagehash"
	"golang.org/x/image/webp"
	"golang.org/x/sync/errgroup"

	"github.com/khodand/memenitpicker_bot/internal/entity"
	"github.com/khodand/memenitpicker_bot/internal/store"
)

type Service struct {
	memes store.Memes
}

func New(memes store.Memes) *Service {
	return &Service{memes: memes}
}

func (s *Service) InsertOrReturnDuplicate(
	ctx context.Context,
	messageID int,
	chatID int64,
	file io.ReadCloser,
) (int, error) {
	img, err := decodeImage(file)
	if err != nil {
		return 0, err
	}

	var dID, aID, pID int
	eg := &errgroup.Group{}
	eg.Go(func() error {
		var err error
		dID, err = s.insert(ctx, messageID, chatID, entity.HashKindDifference, img)
		return err
	})
	eg.Go(func() error {
		var err error
		aID, err = s.insert(ctx, messageID, chatID, entity.HashKindAverage, img)
		return err
	})
	eg.Go(func() error {
		var err error
		pID, err = s.insert(ctx, messageID, chatID, entity.HashKindPerception, img)
		return err
	})
	if err := eg.Wait(); err != nil {
		return 0, err
	}

	return consensus(dID, aID, pID), nil
}

func (s *Service) InsertReturnALl(
	ctx context.Context,
	messageID int,
	chatID int64,
	file io.ReadCloser,
) (int, int, int, error) {
	img, err := decodeImage(file)
	if err != nil {
		return 0, 0, 0, err
	}

	var dID, aID, pID int
	eg := &errgroup.Group{}
	eg.Go(func() error {
		var err error
		dID, err = s.insert(ctx, messageID, chatID, entity.HashKindDifference, img)
		return err
	})
	eg.Go(func() error {
		var err error
		aID, err = s.insert(ctx, messageID, chatID, entity.HashKindAverage, img)
		return err
	})
	eg.Go(func() error {
		var err error
		pID, err = s.insert(ctx, messageID, chatID, entity.HashKindPerception, img)
		return err
	})
	if err := eg.Wait(); err != nil {
		return 0, 0, 0, err
	}

	return dID, aID, pID, nil
}

func (s *Service) insert(
	ctx context.Context,
	messageID int,
	chatID int64,
	hashKind entity.HashKind,
	img image.Image,
) (int, error) {
	var (
		hash      uint64
		imageHash *goimagehash.ImageHash
		err       error
	)

	switch hashKind {
	case entity.HashKindDifference:
		imageHash, err = goimagehash.DifferenceHash(img)
	case entity.HashKindAverage:
		imageHash, err = goimagehash.AverageHash(img)
	case entity.HashKindPerception:
		imageHash, err = goimagehash.PerceptionHash(img)
	default:
		return 0, fmt.Errorf("unknown hash kind: %s", hashKind)
	}
	if err != nil {
		return 0, err
	}
	hash = imageHash.GetHash()

	inserted, err := s.memes.Insert(ctx, &entity.Meme{
		Hash:      hash,
		HashKind:  hashKind,
		ChatID:    chatID,
		MessageID: messageID,
	})
	if err != nil {
		return 0, fmt.Errorf("insert meme: %w", err)
	}

	return inserted.MessageID, nil
}

func decodeImage(file io.ReadCloser) (image.Image, error) {
	img, err := jpeg.Decode(file)
	if err == nil {
		return img, nil
	}
	img, err = png.Decode(file)
	if err == nil {
		return img, nil
	}
	img, _, err = image.Decode(file)
	if err == nil {
		return img, nil
	}
	img, err = webp.Decode(file)
	if err == nil {
		return img, nil
	}

	return img, err
}

func consensus(ids ...int) int {
	minVal := math.MaxInt
	nonZeroCount := 0
	for _, id := range ids {
		if id != 0 {
			nonZeroCount++
			if id < minVal {
				minVal = id
			}
		}
	}
	const consensusThreshold = 2
	if nonZeroCount >= consensusThreshold {
		return minVal
	}
	return 0
}
