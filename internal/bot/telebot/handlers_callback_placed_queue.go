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

func (c *controller) handlePlacedQueueBtnNew() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "placed_queue_btn_new")

		// get queue data
		queueName, queueMemberCount, err := pqBtnNew.parseData(ctx.Callback().Data)
		if err != nil {
			return fmt.Errorf("error parsing placed queue new btn from callback: %w", err)
		}

		// try to add queue
		queueID := queue.GenID(queueName)
		err = c.queueDispatcher.AddPlaced(queueID, queueMemberCount)
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

		// callbackBundle := c.bundle.Callback()

		// send message
		if err := ctx.Edit(c.bundle.Callback().PlacedQueue().Main(queueName), c.placedQueueMarkup(queueID, queueName, make([]telegram.Person, queueMemberCount))); err != nil && !strings.Contains(err.Error(), "True") {
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

func (c *controller) handlePlacedQueueBtnSubmit() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		defer ctx.Respond()
		ctx.Set("handler_type", "callback")
		ctx.Set("handler", "placed_queue_btn_submit")

		// getting queue data
		queueID, queueName, memberPosition, err := pqBtnSubmit.parseData(ctx.Callback().Data)
		if err != nil {
			return fmt.Errorf("error parsing placed queue submit btn from callback: %w", err)
		}

		// submit person and get list
		sender := ctx.Callback().Sender
		list, err := c.queueDispatcher.SubmitPlacedSenderAndList(
			queue.ID(queueID),
			telegram.SenderID(sender.ID),
			telegram.Person{
				Username:  sender.Username,
				FirstName: sender.FirstName,
				LastName:  sender.LastName,
			}, memberPosition)
		if errors.Is(err, dispatcher.ErrQueueMemberCountIncorrect) {
			// errorBundle := c.bundle.Errors()
			// ctx.RespondText(errorBundle.SubmitAgain())
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
			return fmt.Errorf("error submitting placed sender or getting list: %w", err)
		}

		// edit message
		if err := ctx.Edit(c.bundle.Callback().PlacedQueue().Main(queueName), &tele.SendOptions{
			ReplyMarkup:           c.placedQueueMarkup(queue.ID(queueID), queueName, list),
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		}); err != nil && !strings.Contains(err.Error(), "True") {
			return fmt.Errorf("error editing message: %w", err)
		}
		return nil
	}
}
