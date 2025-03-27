package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/subliker/que-bot/internal/bot/telebot"
	"github.com/subliker/que-bot/internal/logger/zap"
	"github.com/subliker/que-bot/internal/validation"
)

// Config is struct that contain app options
type Config struct {
	Logger zap.Config     `yaml:"logger" env-prefix:"QUE_LOGGER_"`
	Bot    telebot.Config `yaml:"bot" env-required:"true" env-prefix:"QUE_BOT_"`
}

var _filePath = flag.String("config", "configs/config.yml", "config file path")
var _help = flag.Bool("help", false, "show configuration help")

// Load reads config from file and env into struct
func Load() Config {
	logger := zap.Logger.WithFields("layer", "logger")
	var cfg Config

	// if help mode
	if *_help {
		// getting help text
		headerText := "Configuration options:"
		help, err := cleanenv.GetDescription(&cfg, &headerText)
		if err != nil {
			logger.Fatalf("error getting configuration description: %s", err.Error())
		}
		// print help text and exit
		fmt.Println(help)
		os.Exit(0)
	}

	// reading config from file
	err := cleanenv.ReadConfig(*_filePath, &cfg)
	if err != nil {
		logger.Fatalf("error loading config: %s", err)
	}

	// reading config from env
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		logger.Fatalf("error loading config: %s", err)
	}

	// validate config
	if err := validation.Validate.Struct(&cfg); err != nil {
		logger.Fatalf("error config validation: %s", err)
	}
	return cfg
}
