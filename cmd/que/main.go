package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/subliker/que-bot/internal/app"
	"github.com/subliker/que-bot/internal/bot/telebot"
	"github.com/subliker/que-bot/internal/config"
	"github.com/subliker/que-bot/internal/dispatcher"
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

	// making queue dispatcher
	qd := dispatcher.NewQueueDispatcher(cfg.Dispatcher)

	// making bot controller
	bc, err := telebot.NewController(logger, cfg.Bot, qd)
	if err != nil {
		logger.Fatalf("error making bot controller: %s", err)
	}

	go func() {
		sm := &http.ServeMux{}
		sm.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "it works!")
		}))
		http.ListenAndServe(":8080", sm)
	}()

	// creating app
	a := app.New(logger, bc)
	// running app
	a.Run(context.Background())

}
