package telebot

// Config contains bot controller configuration.
type Config struct {
	Token             string `yaml:"token" env:"TOKEN" env-required:"true" env-description:"Token is bot father telegram token for bot api"`
	LongPollerTimeout int    `yaml:"long_poller_timeout" env:"LONG_POLLER_TIMEOUT" env-default:"10" env-description:"LongPollerTimeout is time(in seconds) of response waiting"`
}
