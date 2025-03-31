package telebot

import (
	tele "gopkg.in/telebot.v4"
)

func (c *controller) middlewareRecover(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Error("recovery: panic ", r)
			}
		}()
		return next(ctx)
	}
}
