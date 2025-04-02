package telebot

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	tele "gopkg.in/telebot.v4"
)

var queueBtnNew = tele.Btn{
	Text:   "new",
	Unique: "new",
}

func (c *controller) queueBtnNewMarkup(queueName string) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := queueBtnNew
	btn.Text = c.bundle.Query().Btns().New()
	btn = queueBtnNewData(btn, queueName)
	mk.Inline(tele.Row{btn})
	return mk
}

func queueBtnNewData(btn tele.Btn, queueName string) tele.Btn {
	btn.Data = queueName
	return btn
}

func (c *controller) handleQuery() tele.HandlerFunc {
	return func(ctx tele.Context) error {
		ctx.Set("handler_type", "query")
		ctx.Set("handler", "queue_query")

		queryBundle := c.bundle.Query()

		// format queue name from query
		queueName := queryTextRegexp.ReplaceAllString(ctx.Query().Text, "")
		queueName = strings.TrimSpace(queueName)

		// 62 = all
		// 11 = 8 bytes of id + 2 byte of | + 2 byte(now max) of btn unique
		if len(queueName) > 62-11 {
			queueName = queueName[:62-11]
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
		if err := ctx.Answer(&tele.QueryResponse{
			Results: tele.Results{
				&tele.ArticleResult{
					Title:       queryBundle.Main().Title(queueName),
					Description: queryBundle.Main().Description(),
					Text:        queryBundle.Main().Text(queueName),
					ResultBase: tele.ResultBase{
						ReplyMarkup: c.queueBtnNewMarkup(queueName),
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
