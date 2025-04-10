package telebot

import (
	tele "gopkg.in/telebot.v4"
)

// initHandlers initializes all bot controller message handlers
func (c *controller) initHandlers() {
	c.client.Handle("/start", c.handleStart())
	c.client.Handle(tele.OnQuery, c.handleQuery())

	c.client.Handle(&queueBtnSubmit, c.handleQueueBtnSubmit())
	c.client.Handle(&queueBtnRemove, c.handleQueueBtnRemove())
	c.client.Handle(&queueBtnNew, c.handleQueueBtnNew())

	c.client.Handle(&placedQueueBtnNew, c.handlePlacedQueueBtnNew())
}
