package telebot

// initHandlers initializes all bot controller message handlers
func (c *controller) initHandlers() {
	c.client.Handle("/start", c.handleStart())
}
