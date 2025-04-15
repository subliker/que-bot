package telebot

import (
	"fmt"
	"strings"

	"github.com/subliker/que-bot/internal/dispatcher/queue"
	tele "gopkg.in/telebot.v4"
)

// queueBtnNew is btn to create queue with queueName
type queueBtnNew struct {
	btn
}

var qBtnNew = queueBtnNew{
	btn: newBtn("qnw"),
}

func (b *queueBtnNew) setData(queueName string) {
	b.Data = queueName
}

func (b *queueBtnNew) parseData(callbackData string) (queueName string) {
	return callbackData
}

func (c *controller) queueBtnNewMarkup(queueName string) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := qBtnNew
	btn.Text = c.bundle.Query().Btns().New()
	btn.setData(queueName)
	mk.Inline(tele.Row{btn.tele()})
	return mk
}

// queueBtnSubmit is btn to submit user in queue with queue id
type queueBtnSubmit struct {
	btn
}

var qBtnSubmit = queueBtnSubmit{
	btn: newBtn("qst"),
}

func (b *queueBtnSubmit) setData(queueID, queueName string) {
	b.Data = strings.Join([]string{string(queueID), queueName}, "|")
}

func (b *queueBtnSubmit) parseData(callbackData string) (queueID, queueName string, err error) {
	ss := strings.Split(callbackData, "|")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("incorrect callback data: %s", callbackData)
	}

	return ss[0], ss[1], nil
}

func (c *controller) queueSubmitFirstMarkup(queueID queue.ID, queueName string) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := qBtnSubmit
	btn.Text = c.bundle.Callback().Btns().SubmitFirst()
	btn.setData(string(queueID), queueName)
	mk.Inline(tele.Row{btn.tele()})
	return mk
}

// queueBtnRemove is btn to remove user from queue
type queueBtnRemove struct {
	btn
}

var qBtnRemove = queueBtnRemove{
	btn: newBtn("qbr"),
}

func (b *queueBtnRemove) setData(queueID, queueName string) {
	b.Data = strings.Join([]string{string(queueID), queueName}, "|")
}

func (b *queueBtnRemove) parseData(callbackData string) (queueID, queueName string, err error) {
	ss := strings.Split(callbackData, "|")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("incorrect callback data: %s", callbackData)
	}

	return ss[0], ss[1], nil
}

func (c *controller) queueMarkup(queueID queue.ID, queueName string, lenList int) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btnSubmit, btnRemove := qBtnSubmit, qBtnRemove
	btnSubmit.Text = c.bundle.Callback().Btns().Submit(lenList + 1)
	btnRemove.Text = c.bundle.Callback().Btns().Remove()
	btnSubmit.setData(string(queueID), queueName)
	btnRemove.setData(string(queueID), queueName)
	mk.Inline(tele.Row{btnSubmit.tele()}, tele.Row{btnRemove.tele()})
	return mk
}
