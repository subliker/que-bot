package main

import (
	"context"
	"flag"

	"github.com/subliker/que-bot/internal/app"
	"github.com/subliker/que-bot/internal/bot/telebot"
	"github.com/subliker/que-bot/internal/config"
	"github.com/subliker/que-bot/internal/logger/zap"
)

func main() {
	flag.Parse()

	// reading config
	cfg := config.Load()

	// creating logger
	logger := zap.NewLogger(cfg.Logger, "que-bot")
	// update global logger
	zap.Logger = logger

	// making bot controller
	bc, err := telebot.NewController(logger, cfg.Bot)
	if err != nil {
		logger.Fatalf("error making bot controller: %s", err)
	}

	// creating app
	a := app.New(logger, bc)
	// running app
	a.Run(context.Background())

}
