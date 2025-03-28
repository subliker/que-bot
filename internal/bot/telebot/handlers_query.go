package telebot

import (
	"fmt"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v4"
)

func (c *controller) handleQuery() tele.HandlerFunc {
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

		// handle answer
		mk := c.client.NewMarkup()
		btn := queueQueryBtnNew
		// generate uniq queue uuid
		btn.Data = uuid.NewString()
		mk.Inline(tele.Row{btn})
		if err := ctx.Answer(&tele.QueryResponse{
			Results: tele.Results{
				&tele.ArticleResult{
					Title:       fmt.Sprintf("Queue %s", ctx.Query().Text),
					Description: "Create group query",
					Text:        "push button below to create queue",
					ResultBase: tele.ResultBase{
						ReplyMarkup: mk,
					},
				},
			},
			CacheTime: 1,
		}); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

		return nil
	}
}
