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

var queueBtnSubmit = tele.Btn{
	Text:   "st",
	Unique: "st",
}

func queueBtnSubmitData(btn tele.Btn, queueID, queueName string) tele.Btn {
	btn.Data = strings.Join([]string{string(queueID), queueName}, "|")
	return btn
}

func (c *controller) queueSubmitFirstMarkup(queueID queue.ID, queueName string) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := queueBtnSubmit
	btn.Text = c.bundle.Callback().Btns().SubmitFirst()
	btn.Data = strings.Join([]string{string(queueID), queueName}, "|")
	mk.Inline(tele.Row{btn})
	return mk
}

var queueBtnRemove = tele.Btn{
	Text:   "rm",
	Unique: "rm",
}

func queueBtnRemoveData(btn tele.Btn, queueID, queueName string) tele.Btn {
	btn.Data = strings.Join([]string{string(queueID), queueName}, "|")
	return btn
}

func (c *controller) queueMarkup(queueID queue.ID, queueName string, lenList int) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btnSubmit, btnRemove := queueBtnSubmit, queueBtnRemove
	btnSubmit.Text = c.bundle.Callback().Btns().Submit(lenList + 1)
	btnRemove.Text = c.bundle.Callback().Btns().Remove()
	btnSubmit = queueBtnSubmitData(btnSubmit, string(queueID), queueName)
	btnRemove = queueBtnRemoveData(btnRemove, string(queueID), queueName)
	mk.Inline(tele.Row{btnSubmit, btnRemove})
	return mk
}

func (c *controller) queueText(queueName string, list []telegram.Person) string {
	txt := c.bundle.Callback().Queue().Head(queueName) + "\n"
	for i, p := range list {
		txt += c.bundle.Callback().Queue().Member(i+1, p.FirstName, p.LastName, p.Username) + "\n"
	}
	return txt
}

func (c *controller) handleQueueBtnNew() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "queue_btn_new")

		// get queue data
		queueName := ctx.Callback().Data

		// try to add queue
		queueID := queue.GenID(queueName)
		err := c.queueDispatcher.Add(queueID)
		if errors.Is(err, dispatcher.ErrQueueAlreadyExists) {
			errorBundle := c.bundle.Errors()
			ctx.RespondAlert(errorBundle.QueueIdCollision())
			return fmt.Errorf("queue for query id already exists")
		}
		if err != nil {
			errorBundle := c.bundle.Errors()
			ctx.RespondAlert(errorBundle.Internal())
			return fmt.Errorf("error adding queue in dispatcher: %w", err)
		}

		callbackBundle := c.bundle.Callback()

		// send message
		if err := ctx.Edit(callbackBundle.QueueNew().Main(queueName), c.queueSubmitFirstMarkup(queueID, queueName)); err != nil && !strings.Contains(err.Error(), "True") {
			if strings.Contains(err.Error(), "BUTTON_DATA_INVALID") {
				errorsBundle := c.bundle.Errors()
				ctx.Edit(errorsBundle.ButtonDataLength() +
					errorsBundle.Tail())
			}
			return err
		}
		return nil
	}
}

func (c *controller) handleQueueBtnSubmit() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "queue_btn_submit")

		// getting queue data
		data := strings.Split(ctx.Callback().Data, "|")
		if len(data) != 2 {
			return fmt.Errorf("callback length data arguments error")
		}
		queueID, queueName := queue.ID(data[0]), data[1]

		// submit person and get list
		sender := ctx.Callback().Sender
		list, err := c.queueDispatcher.SubmitSenderAndList(
			queueID,
			telegram.SenderID(sender.ID),
			telegram.Person{
				Username:  sender.Username,
				FirstName: sender.FirstName,
				LastName:  sender.LastName,
			})
		if errors.Is(err, dispatcher.ErrQueueSenderAlreadyExists) {
			errorBundle := c.bundle.Errors()
			ctx.RespondText(errorBundle.SubmitAgain())
			return nil
		}
		if errors.Is(err, dispatcher.ErrQueueNotExists) {
			errorBundle := c.bundle.Errors()
			ctx.RespondAlert(errorBundle.QueueNotFound())
			return nil
		}
		if err != nil {
			errorBundle := c.bundle.Errors()
			ctx.RespondAlert(errorBundle.Internal())
			return fmt.Errorf("error submitting sender or getting list: %w", err)
		}

		// edit message
		if err := ctx.Edit(c.queueText(queueName, list), &tele.SendOptions{
			ReplyMarkup:           c.queueMarkup(queueID, queueName, len(list)),
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		}); err != nil && !strings.Contains(err.Error(), "True") {
			return fmt.Errorf("error editing message: %w", err)
		}
		return nil
	}
}

func (c *controller) handleQueueBtnRemove() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "queue_query_btn_remove")

		// getting queue data
		data := strings.Split(ctx.Callback().Data, "|")
		if len(data) != 2 {
			return fmt.Errorf("callback length data arguments error")
		}
		queueID, queueName := queue.ID(data[0]), data[1]

		// submit person and get list
		sender := ctx.Callback().Sender
		list, err := c.queueDispatcher.RemoveSenderAndList(
			queueID,
			telegram.SenderID(sender.ID))
		if errors.Is(err, dispatcher.ErrQueueSenderNotExists) {
			errorBundle := c.bundle.Errors()
			ctx.RespondText(errorBundle.RemoveIfNot())
			return nil
		}
		if errors.Is(err, dispatcher.ErrQueueNotExists) {
			errorBundle := c.bundle.Errors()
			ctx.RespondAlert(errorBundle.QueueNotFound())
			return nil
		}
		if err != nil {
			errorBundle := c.bundle.Errors()
			ctx.RespondAlert(errorBundle.Internal())
			return fmt.Errorf("error removing sender or getting list: %w", err)
		}

		// edit message
		err = ctx.Edit(c.queueText(queueName, list), &tele.SendOptions{
			ReplyMarkup:           c.queueMarkup(queueID, queueName, len(list)),
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "True") {
				return fmt.Errorf("error editing message: %w", err)
			}
		}

		return nil
	}
}
