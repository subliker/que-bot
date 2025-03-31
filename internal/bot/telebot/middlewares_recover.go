package telebot

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func (c *controller) middlewareRecover(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer func() {
			if r := recover(); r != nil {
				c.onError(fmt.Errorf("recovery: panic %s", r), ctx)
			}
		}()
		return next(ctx)
	}
}
