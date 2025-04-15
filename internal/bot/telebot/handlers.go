package telebot

import (
	tele "gopkg.in/telebot.v4"
)

// initHandlers initializes all bot controller message handlers
func (c *controller) initHandlers() {
	c.client.Handle("/start", c.handleStart())
	c.client.Handle(tele.OnQuery, c.handleQuery())

	c.client.Handle(&qBtnSubmit, c.handleQueueBtnSubmit())
	c.client.Handle(&qBtnRemove, c.handleQueueBtnRemove())
	c.client.Handle(&qBtnNew, c.handleQueueBtnNew())

	c.client.Handle(&pqBtnSubmit, c.handlePlacedQueueBtnSubmit())
	c.client.Handle(&pqBtnSubmitHead, c.handlePlacedQueueBtnSubmitHead())
	c.client.Handle(&pqBtnRemove, c.handlePlacedQueueBtnRemove())
	c.client.Handle(&pqBtnNew, c.handlePlacedQueueBtnNew())
}
