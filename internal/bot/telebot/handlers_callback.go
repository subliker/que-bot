package telebot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/dispatcher/queue"
	"github.com/subliker/que-bot/internal/domain/telegram"
	tele "gopkg.in/telebot.v4"
)

var queueQueryBtnNew = tele.Btn{
	Text:   "new",
	Unique: "new",
}

func queueQueryBtnNewData(btn tele.Btn, queueName string) tele.Btn {
	btn.Data = queueName
	return btn
}

func (c *controller) handleQueueQueryBtnNew() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "queue_query_btn_new")

		// get queue data
		queueName := ctx.Callback().Data

		// try to add queue
		queueID := queue.GenID(queueName)
		err := c.queueDispatcher.Add(queueID)
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
		btn.Data = strings.Join([]string{string(queueID), queueName}, "|")
		mk.Inline(tele.Row{btn})
		if err := ctx.Edit(callbackBundle.QueueNew().Main(queueName), mk); err != nil && !strings.Contains(err.Error(), "True") {
			if strings.Contains(err.Error(), "BUTTON_DATA_INVALID") {
				errorsBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Errors()
				ctx.Edit(errorsBundle.ButtonDataLength() +
					errorsBundle.Tail())
			}
			return err
		}
		return nil
	}
}

var queueQueryBtnSubmit = tele.Btn{
	Text:   "submit",
	Unique: "submit",
}

func queueQueryBtnSubmitData(btn tele.Btn, queueID, queueName string) tele.Btn {
	btn.Data = strings.Join([]string{string(queueID), queueName}, "|")
	return btn
}

var queueQueryBtnSubmitLength = len(strings.Join([]string{
	queueQueryBtnSubmit.Unique, string(queue.GenID("")),
}, "|"))

func (c *controller) handleQueueQueryBtnSubmit() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "queue_query_btn_submit")

		// getting queue data
		data := strings.Split(ctx.Callback().Data, "|")
		if len(data) != 2 {
			return fmt.Errorf("callback length data arguments error")
		}
		queueID, queueName := queue.ID(data[0]), data[1]

		callbackBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Callback()

		// submit person and get list
		sender := ctx.Callback().Sender
		lst, err := c.queueDispatcher.SubmitSenderAndList(
			queueID,
			telegram.SenderID(sender.ID),
			telegram.Person{
				Username:  sender.Username,
				FirstName: sender.FirstName,
				LastName:  sender.LastName,
			})
		if errors.Is(err, dispatcher.ErrQueueSenderAlreadyExists) {
			errorBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Errors()
			ctx.RespondText(errorBundle.SubmitAgain())
			return nil
		}
		if errors.Is(err, dispatcher.ErrQueueNotExists) {
			errorBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Errors()
			// make recreate button
			mk := c.client.NewMarkup()
			btn := queueQueryBtnSubmit
			btn.Text = callbackBundle.Btns().Relive()
			btn = queueQueryBtnSubmitData(btn, string(queueID), queueName)
			mk.Inline(tele.Row{btn})
			ctx.RespondAlert(errorBundle.QueueNotFound())
			return nil
		}
		if err != nil {
			errorBundle := c.langBundle(ctx.Callback().Sender.LanguageCode).Errors()
			ctx.RespondAlert(errorBundle.Internal())
			return fmt.Errorf("error submitting sender or getting list: %w", err)
		}

		// format answer
		txt := callbackBundle.Queue().Head(queueName) + "\n"
		for i, p := range lst {
			txt += callbackBundle.Queue().Member(i+1, p.FirstName, p.LastName, p.Username) + "\n"
		}

		// edit message
		mk := c.client.NewMarkup()
		btn := queueQueryBtnSubmit
		btn.Text = callbackBundle.Btns().Submit(len(lst) + 1)
		btn = queueQueryBtnSubmitData(btn, string(queueID), queueName)
		mk.Inline(tele.Row{btn})
		ctx.Edit(txt, &tele.SendOptions{
			ReplyMarkup:           mk,
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})
		return nil
	}
}
