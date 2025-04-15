package telebot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/dispatcher/queue"
	"github.com/subliker/que-bot/internal/domain/telegram"
	tele "gopkg.in/telebot.v4"
)

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
		queueName := qBtnNew.parseData(ctx.Callback().Data)

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
		queueID, queueName, err := qBtnSubmit.parseData(ctx.Callback().Data)
		if err != nil {
			return fmt.Errorf("error parsing queue submit btn from callback: %w", err)
		}

		// submit person and get list
		sender := ctx.Callback().Sender
		list, err := c.queueDispatcher.SubmitSenderAndList(
			queue.ID(queueID),
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

		// edit message with limiter
		c.limiter.Do(string(queueID), func() {
			if err := ctx.Edit(c.queueText(queueName, list), &tele.SendOptions{
				ReplyMarkup:           c.queueMarkup(queue.ID(queueID), queueName, len(list)),
				DisableWebPagePreview: true,
				ParseMode:             tele.ModeMarkdown,
			}); err != nil && !strings.Contains(err.Error(), "True") {
				if strings.Contains(err.Error(), "retry after") {
					errorBundle := c.bundle.Errors()
					ctx.RespondAlert(errorBundle.RetryAfter())
					return
				}
				c.logger.
					WithFields("handler_type", "callback",
						"handler", "queue_btn_submit").
					Errorf("error editing message: %s", err)
				return
			}
		}, time.Millisecond*1500)

		return nil
	}
}

func (c *controller) handleQueueBtnRemove() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "queue_query_btn_remove")

		// getting queue data
		queueID, queueName, err := qBtnRemove.parseData(ctx.Callback().Data)
		if err != nil {
			return fmt.Errorf("error parsing queue remove btn from callback: %w", err)
		}

		// submit person and get list
		sender := ctx.Callback().Sender
		list, err := c.queueDispatcher.RemoveSenderAndList(
			queue.ID(queueID),
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

		// edit message with limiter
		c.limiter.Do(string(queueID), func() {
			err := ctx.Edit(c.queueText(queueName, list), &tele.SendOptions{
				ReplyMarkup:           c.queueMarkup(queue.ID(queueID), queueName, len(list)),
				DisableWebPagePreview: true,
				ParseMode:             tele.ModeMarkdown,
			})
			if err != nil && !strings.Contains(err.Error(), "True") {
				if strings.Contains(err.Error(), "retry after") {
					errorBundle := c.bundle.Errors()
					ctx.RespondAlert(errorBundle.RetryAfter())
					return
				}
				c.logger.
					WithFields("handler_type", "callback",
						"handler", "queue_btn_submit").
					Errorf("error editing message: %s", err)
				return
			}
		}, time.Millisecond*1500)

		return nil
	}
}
