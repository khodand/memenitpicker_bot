package bot

import (
	"context"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"github.com/khodand/memenitpicker_bot/internal/meme"
	"github.com/khodand/memenitpicker_bot/pkg/telegram"
)

type Service struct {
	client *tele.Bot
	memes  *meme.Service
}

func New(client *tele.Bot, memes *meme.Service) *Service {
	return &Service{
		client: client,
		memes:  memes,
	}
}

func (s *Service) RegisterRoutes() {
	s.client.Handle(tele.OnPhoto, func(c tele.Context) error {
		ctx := context.Background()
		err := s.onPhoto(ctx, c.Message())
		if err != nil {
			telegram.GetLogger(c).Error("on photo", zap.Error(err))
		}
		return nil
	})
}

func (s *Service) onPhoto(ctx context.Context, message *tele.Message) error {
	photo := message.Photo
	if photo == nil {
		return nil
	}
	file, err := s.client.File(&photo.File)
	if err != nil {
		return err
	}

	duplicationID, err := s.memes.InsertOrReturnDuplicate(ctx, message.ID, message.Chat.ID, file)
	if err != nil {
		return err
	}
	if duplicationID == message.ID {
		return nil
	}

	_, err = s.client.Reply(
		&tele.Message{
			ID:   duplicationID,
			Chat: message.Chat,
		},
		"Кажется, этот мем уже присылали)",
	)
	if err != nil {
		return err
	}

	return nil
}
