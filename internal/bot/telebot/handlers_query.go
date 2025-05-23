package telebot

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	tele "gopkg.in/telebot.v4"
)

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
		placedQueueText := queryBundle.IncorrectCount()
		if len(queueSplit) > 0 {
			var err error
			queueMemberCount, err = strconv.Atoi(queueSplit[len(queueSplit)-1])
			if err == nil && queueMemberCount > 0 && queueMemberCount < 99 {
				// if queue member count is set
				queueNameWithMembersCount = strings.Join(queueSplit[:len(queueSplit)-1], " ")
				placedQueueText = queryBundle.PlacedQueue().Text(queueNameWithMembersCount, queueMemberCount)
				placedQueueReplyMarkup = tele.ResultBase{
					ReplyMarkup: c.placedQueueBtnNewMarkup(queueNameWithMembersCount, queueMemberCount),
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
					Text:        placedQueueText,
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
