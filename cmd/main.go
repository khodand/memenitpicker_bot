package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/khodand/memenitpicker_bot/cmd/config"
	"github.com/khodand/memenitpicker_bot/internal/bot"
	"github.com/khodand/memenitpicker_bot/internal/meme"
	"github.com/khodand/memenitpicker_bot/internal/store"
	"github.com/khodand/memenitpicker_bot/pkg/logger"
	"github.com/khodand/memenitpicker_bot/pkg/postgres"
	ptelegram "github.com/khodand/memenitpicker_bot/pkg/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(fmt.Errorf("failed to init config: %w", err))
	}

	log := logger.New(cfg.General.Debug)

	database, err := postgres.NewPgxPool(cfg.Postgres)
	if err != nil {
		log.Fatal("failed to init database", zap.Error(err))
	}
	memesDB := store.NewMemes(database)

	telegramClient, err := ptelegram.NewClient(log, cfg.Telegram)
	if err != nil {
		log.Fatal("failed to init bot client", zap.Error(err))
	}

	memeService := meme.New(memesDB)
	telegramBot := bot.New(telegramClient, memeService)
	telegramBot.RegisterRoutes()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("bot client is starting...")
		telegramClient.Start()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	telegramClient.Stop()
	wg.Wait()

	log.Info("bot was stopped gracefully")
	_ = log.Sync()
}
