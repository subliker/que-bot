package telebot

import (
	"errors"
	"fmt"
	"sync"

	"github.com/subliker/que-bot/internal/dispatcher"
	"github.com/subliker/que-bot/internal/domain/telegram"
	tele "gopkg.in/telebot.v4"
)

var queueQueryBtnSubmitQueue tele.Btn

func (c *controller) handleQuery() tele.HandlerFunc {
	// make markup
	var once sync.Once
	mk := c.client.NewMarkup()
	once.Do(func() {
		queueQueryBtnSubmitQueue = mk.Data("submit", "submit-queue")
		mk.Inline(mk.Row(queueQueryBtnSubmitQueue))
	})
	// logger := c.logger.WithFields("handler", "handler")
	return func(ctx tele.Context) error {
		// logger := logger.WithFields(
		// 	"sender_id", ctx.Sender().ID,
		// )

		// send error if not in group
		if ctx.Query().ChatType != "group" {
			if err := ctx.Answer(&tele.QueryResponse{
				Results: tele.Results{
					&tele.ArticleResult{
						Title:       fmt.Sprintf("Queue %s", ctx.Query().Text),
						Description: "Create group query",
						Text:        "Queue can be created only in groups",
					},
				},
			}); err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}
			return nil
		}

		// try to add queue
		err := c.queueDispatcher.Add(telegram.QueryID(ctx.Query().ID))
		if errors.Is(err, dispatcher.ErrQueueAlreadyExists) {
			// send error that queue is already exists
			if err := ctx.Answer(&tele.QueryResponse{
				Results: tele.Results{
					&tele.ArticleResult{
						Title:       fmt.Sprintf("Queue %s", ctx.Query().Text),
						Description: "Create group query",
						Text:        "Queue in this chat already exists, you need to finish it before new",
					},
				},
			}); err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}
		}
		if err != nil {
			if err := ctx.Answer(&tele.QueryResponse{
				Results: tele.Results{
					&tele.ArticleResult{
						Title:       fmt.Sprintf("Queue %s", ctx.Query().Text),
						Description: "Create group query",
						Text:        "internal error",
					},
				},
			}); err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}
			return fmt.Errorf("error adding queue in dispatcher: %w", err)
		}

		// handle answer
		if err := ctx.Answer(&tele.QueryResponse{
			Results: tele.Results{
				&tele.ArticleResult{
					Title:       fmt.Sprintf("Queue %s", ctx.Query().Text),
					Description: "Create group query",
					Text:        "push button below to submit queue",
					ResultBase: tele.ResultBase{
						ReplyMarkup: mk,
					},
				},
			},
		}); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

		return nil
	}
}
