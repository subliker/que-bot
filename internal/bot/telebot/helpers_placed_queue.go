package telebot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/subliker/que-bot/internal/dispatcher/queue"
	"github.com/subliker/que-bot/internal/domain/telegram"
	tele "gopkg.in/telebot.v4"
)

// placedQueueBtnNew is btn to create new placed queue
type placedQueueBtnNew struct {
	btn
}

var pqBtnNew = placedQueueBtnNew{
	btn: newBtn("pnw"),
}

func (b *placedQueueBtnNew) setData(queueName string, queueMemberCount int) {
	b.Data = strings.Join([]string{queueName, strconv.Itoa(queueMemberCount)}, "|")
}

func (b *placedQueueBtnNew) parseData(callbackData string) (queueName string, queueMemberCount int, err error) {
	ss := strings.Split(callbackData, "|")
	if len(ss) != 2 {
		return "", 0, fmt.Errorf("incorrect callback data: %s", callbackData)
	}

	n, err := strconv.Atoi(ss[1])
	if err != nil {
		return "", 0, fmt.Errorf("incorrect callback member count: %s", ss[1])
	}

	return ss[0], n, nil
}

func (c *controller) placedQueueBtnNewMarkup(queueName string, queueMemberCount int) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()
	btn := pqBtnNew
	btn.Text = c.bundle.Query().Btns().New()
	btn.setData(queueName, queueMemberCount)
	mk.Inline(tele.Row{btn.tele()})
	return mk
}

// placedQueueBtnSubmitHead is btn to submit from head in queue
type placedQueueBtnSubmitHead struct {
	btn
}

var pqBtnSubmitHead = placedQueueBtnSubmitHead{
	btn: newBtn("psh"),
}

func (b *placedQueueBtnSubmitHead) setData(queueID, queueName string) {
	b.Data = strings.Join([]string{queueID, queueName}, "|")
}

func (b *placedQueueBtnSubmitHead) parseData(callbackData string) (queueID, queueName string, err error) {
	ss := strings.Split(callbackData, "|")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("incorrect callback data: %s", callbackData)
	}

	return ss[0], ss[1], nil
}

// placedQueueBtnSubmit is btn with place to submit user
type placedQueueBtnSubmit struct {
	btn
}

var pqBtnSubmit = placedQueueBtnSubmit{
	btn: newBtn("pst"),
}

func (b *placedQueueBtnSubmit) setData(queueID, queueName string, memberPosition int) {
	b.Data = strings.Join([]string{queueID, queueName, strconv.Itoa(memberPosition)}, "|")
}

func (b *placedQueueBtnSubmit) parseData(callbackData string) (queueID, queueName string, memberPosition int, err error) {
	ss := strings.Split(callbackData, "|")
	if len(ss) != 3 {
		return "", "", 0, fmt.Errorf("incorrect callback data: %s", callbackData)
	}

	n, err := strconv.Atoi(ss[2])
	if err != nil {
		return "", "", 0, fmt.Errorf("incorrect callback member position: %s", ss[2])
	}

	return ss[0], ss[1], n, nil
}

// placedQueueBtnRemove is btn to remove person from places queue
type placedQueueBtnRemove struct {
	btn
}

var pqBtnRemove = placedQueueBtnRemove{
	btn: newBtn("prm"),
}

func (b *placedQueueBtnRemove) setData(queueID, queueName string) {
	b.Data = strings.Join([]string{queueID, queueName}, "|")
}

func (b *placedQueueBtnRemove) parseData(callbackData string) (queueID, queueName string, err error) {
	ss := strings.Split(callbackData, "|")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("incorrect callback data: %s", callbackData)
	}

	return ss[0], ss[1], nil
}

func (c *controller) placedQueueMarkup(queueID queue.ID, queueName string, list []telegram.Person) *tele.ReplyMarkup {
	mk := c.client.NewMarkup()

	columns := 1
	if len(list) > 10 {
		columns = 2
	}

	rs := make([]tele.Row, 1)
	r := tele.Row{}

	i := 0
	placeBtn := func(b tele.Btn) {
		r = append(r, b)
		if len(r) == columns || i == len(list)-1 {
			rs = append(rs, r)
			r = tele.Row{}
		}
	}

	placedCount := 0
	p := telegram.Person{}
	for i, p = range list {
		if p.Username != "" {
			placedCount++
			placeBtn(mk.URL(c.bundle.Callback().PlacedQueue().Member(i+1, p.FirstName, p.LastName), "https://t.me/"+p.Username))
			continue
		}
		submitBtn := pqBtnSubmit
		submitBtn.Text = strconv.Itoa(i+1) + "."
		submitBtn.setData(string(queueID), queueName, i)
		placeBtn(submitBtn.tele())
	}

	submitHeadBtn := pqBtnSubmitHead
	submitHeadBtn.Text = c.bundle.Callback().Btns().Submit(placedCount)
	submitHeadBtn.setData(string(queueID), queueName)
	rs[0] = tele.Row{submitHeadBtn.tele()}

	removeBtn := pqBtnRemove
	removeBtn.Text = c.bundle.Callback().Btns().Remove()
	removeBtn.setData(string(queueID), queueName)
	rs = append(rs, tele.Row{removeBtn.tele()})
	mk.Inline(rs...)
	return mk
}
