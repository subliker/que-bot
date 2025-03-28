package telebot

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/dispatcher/queue"
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
		queueID := ctx.Callback().Data
		if queueID == "" {
			return fmt.Errorf("empty queue query btn submit data")
		}

		// submit person
		uuid, err := uuid.Parse(queueID)
		if err != nil {
			return fmt.Errorf("error parsing queue uuid from callback data: %w", err)
		}
		sender := ctx.Callback().Sender
		err = c.queueDispatcher.SubmitSender(
			queue.ID(uuid),
			telegram.SenderID(sender.ID),
			telegram.Person{
				Username:  sender.Username,
				FirstName: sender.FirstName,
				LastName:  sender.LastName,
			})
		if errors.Is(err, dispatcher.ErrQueueSenderAlreadyExists) {
			return nil
		}
		if errors.Is(err, dispatcher.ErrQueueNotExists) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error submitting sender: %w", err)
		}

		// get updated array
		lst, err := c.queueDispatcher.List(queue.ID(uuid))
		if errors.Is(err, dispatcher.ErrQueueNotExists) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error getting queue list: %w", err)
		}

		// format answer. TODO move to message package
		txt := "Очередь:\n"
		for i, p := range lst {
			txt += fmt.Sprintf("%d. [%s %s](https://t.me/%s)\n", i+1, p.FirstName, p.LastName, p.Username)
		}

		// edit message
		mk := c.client.NewMarkup()
		btn := queueQueryBtnSubmit
		// set queue uuid
		btn.Data = ctx.Callback().Data
		mk.Inline(tele.Row{btn})
		ctx.Edit(txt, &tele.SendOptions{
			ReplyMarkup:           mk,
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})

		// if err := ctx.Send(fmt.Sprintf("%d. @%s", num, ctx.Callback().Sender.Username)); err != nil {
		// 	return fmt.Errorf("error sending message: %w", err)
		// }
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
		uuid, err := uuid.Parse(ctx.Callback().Data)
		if err != nil {
			return fmt.Errorf("error parsing queue uuid from callback data: %w", err)
		}
		err = c.queueDispatcher.Add(queue.ID(uuid))
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
