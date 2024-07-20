package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/khodand/memenitpicker_bot/internal/meme"
	"github.com/khodand/memenitpicker_bot/internal/store"
	"github.com/khodand/memenitpicker_bot/pkg/postgres"
)

func main() {
	cfg, err := Init()
	if err != nil {
		panic(err)
	}

	dbPool, err := postgres.NewPgxPool(cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	memesDB := store.NewMemes(dbPool)
	memeService := meme.New(memesDB)

	err = processFiles(cfg, memeService)
	if err != nil {
		log.Fatal("Failed to process files: ", err)
	} else {
		log.Println("Data has been successfully processed")
	}
}

func processFiles(cfg Config, memeService *meme.Service) error {
	file, err := os.Open(cfg.PrefixPath + cfg.ResultJSON)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var in tInput
	if err = json.Unmarshal(byteValue, &in); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	processed := 0
	ctx := context.Background()
	for _, msg := range in.Messages {
		if msg.Photo != "" {
			processed++
			if err = processImage(ctx, cfg, msg, memeService); err != nil {
				return fmt.Errorf("process image: %w", err)
			}
		}
		if processed%100 == 0 {
			log.Printf("Processed %d messages\n", processed)
		}
	}
	return nil
}

func processImage(ctx context.Context, cfg Config, msg tMessage, memeService *meme.Service) error {
	imageFile, err := os.Open(cfg.PrefixPath + msg.Photo)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer imageFile.Close()

	if _, err = memeService.InsertOrReturnDuplicate(ctx, msg.ID, cfg.ChatID, imageFile); err != nil {
		return fmt.Errorf("insert meme: %w", err)
	}
	return nil
}
