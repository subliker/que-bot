package telebot

import (
	"context"
	"fmt"
	"time"

	"github.com/subliker/logger"
	"github.com/subliker/que-bot/internal/bot"
	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/lang"
	"github.com/subliker/que-bot/internal/limiter"
	tele "gopkg.in/telebot.v4"
)

type controller struct {
	// client interacts with telegram api
	client          *tele.Bot
	queueDispatcher dispatcher.QueueDispatcher
	limiter         limiter.Limiter

	bundle lang.Messages

	logger logger.Logger
}

// NewController creates new bot controller instance
func NewController(logger logger.Logger,
	cfg Config,
	qd dispatcher.QueueDispatcher,
	limiter limiter.Limiter) (bot.Controller, error) {
	var c controller

	// set logger
	c.logger = logger.WithFields("layer", "bot_controller")

	// making telegram bot api client
	client, err := tele.NewBot(tele.Settings{
		Token: cfg.Token,
		Poller: &tele.LongPoller{
			Timeout: time.Second * time.Duration(cfg.LongPollerTimeout),
			AllowedUpdates: []string{
				"message",
				"inline_query",
				"chosen_inline_result",
				"callback_query",
			},
		},
		ParseMode: tele.ModeMarkdown,
		OnError:   c.onError,
	})
	if err != nil {
		return nil, fmt.Errorf("error making telegram bot api client: %w", err)
	}
	c.client = client

	// set middlewares
	c.client.Use(c.middlewareRecover)
	if cfg.Debug {
		c.client.Use(c.middlewareDebug)
	}

	// set queue dispatcher
	c.queueDispatcher = qd

	// set limiter
	c.limiter = limiter

	// set lang bundle
	var ok bool
	c.bundle, ok = lang.MessagesFor(cfg.Lang)
	if !ok {
		return nil, fmt.Errorf("error incorrect language code: %s", cfg.Lang)
	}

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
	c.logger.Info("bot stopped")
}

func (c *controller) onError(err error, ctx tele.Context) {
	logger := c.logger
	if ctx.Get("handler") != nil && ctx.Get("handler_type") != nil {
		handler := ctx.Get("handler")
		handlerType := ctx.Get("handler_type")
		logger = logger.WithFields("handler", handler, "handler_type", handlerType)
	}
	logger.Error(err)
}
