package telebot

import (
	"errors"
	"fmt"
	"strings"

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

		// getting queue data
		data := strings.Split(ctx.Callback().Data, "|")
		if len(data) != 2 {
			return fmt.Errorf("callback length data arguments error")
		}
		queueID := data[0]
		if queueID == "" {
			return fmt.Errorf("empty queue query btn submit data")
		}
		queueName := data[1]

		// submit person and get list
		uuid, err := uuid.Parse(queueID)
		if err != nil {
			return fmt.Errorf("error parsing queue uuid from callback data: %w", err)
		}
		sender := ctx.Callback().Sender
		lst, err := c.queueDispatcher.SubmitSenderAndList(
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
			return fmt.Errorf("error submitting sender or getting list: %w", err)
		}

		callbackBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Callback()

		// format answer
		txt := callbackBundle.Queue().Head(queueName) + "\n"
		for i, p := range lst {
			txt += callbackBundle.Queue().Member(i+1, p.FirstName, p.LastName, p.Username) + "\n"
		}

		// edit message
		mk := c.client.NewMarkup()
		btn := queueQueryBtnSubmit
		btn.Text = callbackBundle.Btns().Submit(len(lst) + 1)
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

		// get queue data
		data := strings.Split(ctx.Callback().Data, "|")
		if len(data) != 2 {
			return fmt.Errorf("callback length data arguments error")
		}
		uuid, err := uuid.Parse(data[0])
		if err != nil {
			return fmt.Errorf("error parsing queue uuid from callback data: %w", err)
		}
		queueName := data[1]

		// try to add queue
		err = c.queueDispatcher.Add(queue.ID(uuid))
		if errors.Is(err, dispatcher.ErrQueueAlreadyExists) {
			return fmt.Errorf("queue for query id already exists")
		}
		if err != nil {
			return fmt.Errorf("error adding queue in dispatcher: %w", err)
		}

		callbackBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Callback()

		// send message
		mk := c.client.NewMarkup()
		btn := queueQueryBtnSubmit
		btn.Text = callbackBundle.Btns().SubmitFirst()
		btn.Data = ctx.Callback().Data
		mk.Inline(tele.Row{btn})
		ctx.Edit(callbackBundle.QueueNew().Main(queueName), mk)
		return nil
	}
}
