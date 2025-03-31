package telebot

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	tele "gopkg.in/telebot.v4"
)

func (c *controller) handleQuery() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		ctx.Set("handler_type", "query")
		ctx.Set("handler", "queue_query")

		queryBundle := c.langBundle(ctx.Query().Sender.LanguageCode).Query()

		// format queue name from query
		queueName := inputText.ReplaceAllString(ctx.Query().Text, "")
		queueName = strings.TrimSpace(queueName)

		if len(queueName) > 62-queueQueryBtnSubmitLength {
			queueName = queueName[:62-queueQueryBtnSubmitLength]
			if !utf8.ValidString(queueName) {
				queueNameRunes := []rune(queueName)
				queueName = string(queueNameRunes[:len(queueNameRunes)-1])
			}
		}

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
		btn.Data = queueName
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
