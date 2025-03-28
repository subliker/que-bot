package telebot

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func (c *controller) handleStart() tele.HandlerFunc {
	logger := c.logger.WithFields("handler", "start")
	return func(ctx tele.Context) error {
		logger := logger.WithFields(
			"sender_id", ctx.Sender().ID,
			"message_id", ctx.Message().ID,
		)

		// send start message
		if err := ctx.Send("add me to group"); err != nil {
			errMsg := fmt.Errorf("error sending message: %w", err)
			logger.Error(errMsg)
			return errMsg
		}

		return nil
	}
}
