package telebot

import tele "gopkg.in/telebot.v4"

func (c *controller) handlePlacedQueueBtnNew() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()

		return nil
	}
}
