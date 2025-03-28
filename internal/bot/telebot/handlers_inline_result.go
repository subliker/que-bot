package telebot

import (
	"errors"
	"fmt"

	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/domain/telegram"
	tele "gopkg.in/telebot.v4"
)

func (c *controller) handleInlineResult() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		// try to add queue
		err := c.queueDispatcher.Add(telegram.QueryID(ctx.InlineResult().ResultID))
		if errors.Is(err, dispatcher.ErrQueueAlreadyExists) {
			return fmt.Errorf("queue for query id already exists")
		}
		if err != nil {
			return fmt.Errorf("error adding queue in dispatcher: %w", err)
		}
		return nil
	}
}
