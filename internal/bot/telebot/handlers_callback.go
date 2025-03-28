package telebot

import (
	"errors"
	"fmt"

	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/domain/telegram"
	tele "gopkg.in/telebot.v4"
)

var queueQueryBtnSubmit = tele.Btn{
	Text:   "submit",
	Unique: "submit-queue",
}

func (c *controller) handleQueueQueryBtnSubmit() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()

		// getting query id
		queryID := ctx.Callback().Data
		if queryID == "" {
			return fmt.Errorf("empty queue query btn submit data")
		}

		num, err := c.queueDispatcher.SubmitSender(telegram.QueryID(queryID), telegram.SenderID(ctx.Callback().Sender.ID))
		if errors.Is(err, dispatcher.ErrQueueSenderAlreadyExists) {
			return nil
		}
		if errors.Is(err, dispatcher.ErrQueueNotExists) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error submitting sender: %w", err)
		}

		if err := ctx.Send(fmt.Sprintf("%d. @%s", num, ctx.Callback().Sender.Username)); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}
		return nil
	}
}

var queueQueryBtnNew = tele.Btn{
	Text:   "new",
	Unique: "new-queue",
}

func (c *controller) handleQueueQueryBtnNew() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()

		// try to add queue
		err := c.queueDispatcher.Add(telegram.QueryID(ctx.Callback().Data))
		if errors.Is(err, dispatcher.ErrQueueAlreadyExists) {
			return fmt.Errorf("queue for query id already exists")
		}
		if err != nil {
			return fmt.Errorf("error adding queue in dispatcher: %w", err)
		}

		// send message
		mk := c.client.NewMarkup()
		btn := queueQueryBtnSubmit
		btn.Data = ctx.Callback().Data
		mk.Inline(tele.Row{btn})
		ctx.Edit("push commit to take a place", mk)
		return nil
	}
}
