package main

import (
	"context"
	"flag"

	"github.com/subliker/logger/zap"
	"github.com/subliker/que-bot/internal/app"
	"github.com/subliker/que-bot/internal/bot/telebot"
	"github.com/subliker/que-bot/internal/config"
	"github.com/subliker/que-bot/internal/dispatcher"
)

func main() {
	flag.Parse()

	// reading config
	cfg := config.Load()

	// creating logger
	logger := zap.NewLogger(cfg.Logger, "que-bot")
	// update global logger
	zap.Logger = logger

	// making queue dispatcher
	qd := dispatcher.NewQueueDispatcher(cfg.Dispatcher, logger)

	// making bot controller
	bc, err := telebot.NewController(logger, cfg.Bot, qd)
	if err != nil {
		logger.Fatalf("error making bot controller: %s", err)
	}

	// creating app
	a := app.New(logger, bc)
	// running app
	a.Run(context.Background())

}
