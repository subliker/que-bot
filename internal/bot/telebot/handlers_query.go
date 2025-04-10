package telebot

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	tele "gopkg.in/telebot.v4"
)

var (
	queueBtnNew = tele.Btn{
		Unique: "new",
	}

	placedQueueBtnNew = tele.Btn{
		Unique: "pnew",
	}
)

func (c *controller) queueBtnNewMarkup(queueName string) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := queueBtnNew
	btn.Text = c.bundle.Query().Btns().New()
	btn = queueBtnNewData(btn, queueName)
	mk.Inline(tele.Row{btn})
	return mk
}

func (c *controller) placedQueueBtnNewMarkup(queueName string) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := placedQueueBtnNew
	btn.Text = c.bundle.Query().Btns().New()
	btn = placedQueueBtnNewData(btn, queueName)
	mk.Inline(tele.Row{btn})
	return mk
}

func queueBtnNewData(btn tele.Btn, queueName string) tele.Btn {
	btn.Data = queueName
	return btn
}

func placedQueueBtnNewData(btn tele.Btn, queueName string) tele.Btn {
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

		// check potential count members
		queueMemberCount := 0
		queueNameWithMembersCount := ""
		queueSplit := strings.Split(queueName, " ")
		var placedQueueReplyMarkup tele.ResultBase
		if len(queueSplit) > 0 {
			var err error
			queueMemberCount, err = strconv.Atoi(queueSplit[len(queueSplit)-1])
			if err == nil {
				// if queue member count is set
				queueNameWithMembersCount = strings.Join(queueSplit[:len(queueSplit)-1], " ")
				placedQueueReplyMarkup = tele.ResultBase{
					ReplyMarkup: c.placedQueueBtnNewMarkup(queueName),
				}
			}
		}

		// send error if not in group
		if !slices.Contains([]string{
			"group", "supergroup", "channel",
		}, ctx.Query().ChatType) {
			if err := ctx.Answer(&tele.QueryResponse{
				Results: tele.Results{
					&tele.ArticleResult{
						Title:       queryBundle.Queue().Title(queueName),
						Description: queryBundle.Queue().Description(),
						Text:        queryBundle.TextNoGroup(),
					},
					&tele.ArticleResult{
						Title:       queryBundle.PlacedQueue().Title(queueNameWithMembersCount, queueMemberCount),
						Description: queryBundle.PlacedQueue().Description(),
						Text:        queryBundle.TextNoGroup(),
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
					Title:       queryBundle.Queue().Title(queueName),
					Description: queryBundle.Queue().Description(),
					Text:        queryBundle.Queue().Text(queueName),
					ResultBase: tele.ResultBase{
						ReplyMarkup: c.queueBtnNewMarkup(queueName),
					},
				},
				&tele.ArticleResult{
					Title:       queryBundle.PlacedQueue().Title(queueNameWithMembersCount, queueMemberCount),
					Description: queryBundle.PlacedQueue().Description(),
					Text:        queryBundle.PlacedQueue().Text(queueNameWithMembersCount, queueMemberCount),
					ResultBase:  placedQueueReplyMarkup,
				},
			},
			CacheTime: 1,
		}); err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

		return nil
	}
}
