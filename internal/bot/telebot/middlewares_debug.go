package telebot

import (
	"github.com/kr/pretty"
	tele "gopkg.in/telebot.v4"
)

func (c *controller) middlewareDebug(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		c.logger.Debugf("%# v\n", pretty.Formatter(ctx.Update()))
		return next(ctx)
	}
}
