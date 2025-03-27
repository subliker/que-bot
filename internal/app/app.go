package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/subliker/que-bot/internal/bot"
	"github.com/subliker/que-bot/internal/logger"
)

// App is interface to interact with application layers.
// It connects, starts and graceful shutdown the application layers.
type App interface {
	// Run starts app and stops all application layers when ctx done.
	Run(context.Context)
}

type app struct {
	botController bot.Controller
	logger        logger.Logger
}

// New creates new instance of app
func New(logger logger.Logger,
	botController bot.Controller) App {
	var a app

	// set logger
	a.logger = logger.WithFields("layer", "app")

	// set bot controller
	a.botController = botController

	return &a
}

func (a *app) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// receive sys signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Run bot controller
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.botController.Run(ctx)
		cancel()
	}()

	a.logger.Info("app running...")

	// wait until signal will come or context will end
	select {
	case <-quit:
		a.logger.Info("shutdown signal received")
	case <-ctx.Done():
		a.logger.Info("context canceled")
	}

	a.logger.Info("stopping all services...")
	cancel()

	// wait until services will be stopped
	wg.Wait()

	a.logger.Info("app was successfully shutdown :)")
}
