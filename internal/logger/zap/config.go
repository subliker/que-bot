package zap

// Config is struct to configure logger
type Config struct {
	Targets []string `yaml:"targets" env:"TARGETS" env-default:"" env-description:"Targets(addresses) to which the logger sends logs"`
	Dir     string   `yaml:"dir" env:"DIR" env-default:"logs" env-description:"Directory to store logs"`
	Debug   bool     `yaml:"debug" env:"DEBUG" env-default:"false" env-description:"Debug mode"`
}
