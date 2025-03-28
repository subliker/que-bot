package telebot

import (
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v4"
)

func (c *controller) handleQuery() tele.HandlerFunc {
	// logger := c.logger.WithFields("handler", "handler")
	return func(ctx tele.Context) error {
		// logger := logger.WithFields(
		// 	"sender_id", ctx.Sender().ID,
		// )

		queryBundle := c.langBundle(ctx.Query().Sender.LanguageCode).Query()
		queueName := inputText.ReplaceAllString(ctx.Query().Text, "")

		// send error if not in group
		if !slices.Contains([]string{
			"group", "supergroup", "channel",
		}, ctx.Query().ChatType) {
			if err := ctx.Answer(&tele.QueryResponse{
				Results: tele.Results{
					&tele.ArticleResult{
						Title:       queryBundle.Main().Title(queueName),
						Description: queryBundle.Main().Description(),
						Text:        queryBundle.Main().TextNoGroup(),
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
		btn.Text = queryBundle.Btns().New()
		// generate uniq queue uuid
		btn.Data = strings.Join([]string{uuid.NewString(), queueName}, "|")
		mk.Inline(tele.Row{btn})
		if err := ctx.Answer(&tele.QueryResponse{
			Results: tele.Results{
				&tele.ArticleResult{
					Title:       queryBundle.Main().Title(queueName),
					Description: queryBundle.Main().Description(),
					Text:        queryBundle.Main().Text(queueName),
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
