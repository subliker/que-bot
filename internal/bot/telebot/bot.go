package telebot

import (
	"context"
	"fmt"
	"time"

	"github.com/subliker/que-bot/internal/bot"
	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/logger"
	tele "gopkg.in/telebot.v4"
)

type controller struct {
	// client interacts with telegram api
	client          *tele.Bot
	queueDispatcher dispatcher.QueueDispatcher

	logger logger.Logger
}

// NewController creates new bot controller instance
func NewController(logger logger.Logger, cfg Config, qd dispatcher.QueueDispatcher) (bot.Controller, error) {
	var c controller

	// set logger
	c.logger = logger.WithFields("layer", "bot_controller")

	// making telegram bot api client
	client, err := tele.NewBot(tele.Settings{
		Token: cfg.Token,
		Poller: &tele.LongPoller{
			Timeout:        time.Second * time.Duration(cfg.LongPollerTimeout),
			AllowedUpdates: []string{"inline_query"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error making telegram bot api client: %w", err)
	}
	c.client = client

	// set queue dispatcher
	c.queueDispatcher = qd

	// initialize handlers
	c.initHandlers()

	c.logger.Info("bot controller and handlers were initialized")
	return &c, nil
}

func (c *controller) Run(ctx context.Context) {
	c.logger.Info("starting bot...")
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c.client.Start()
		cancel()
	}()
	<-ctx.Done()
}
