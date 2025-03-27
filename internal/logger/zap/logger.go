package zap

import (
	"github.com/subliker/que-bot/internal/logger"
	"go.uber.org/zap"
)

// Logger is global logger with default configuration
var Logger logger.Logger

func init() {
	Logger = NewLogger(Config{}, "")
}

type zapLogger struct {
	logger *zap.SugaredLogger
}
