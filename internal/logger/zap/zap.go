package zap

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/subliker/que-bot/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates sugared zap logger with common config.
// It logs into writer from params.
func NewLogger(cfg Config, serviceName string) logger.Logger {
	var logFile *os.File
	if cfg.Dir != "" {
		// making log file
		os.MkdirAll(cfg.Dir, os.ModePerm)

		var err error
		logFile, err = os.OpenFile(filepath.Join(cfg.Dir, "main.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("error opening log file(%s): %s", logFile.Name(), err)
		}
	}

	// making encoder config
	var zcfg zapcore.EncoderConfig
	var level zapcore.Level
	if cfg.Debug {
		zcfg = zap.NewDevelopmentEncoderConfig()
		level = zapcore.DebugLevel
	} else {
		zcfg = zap.NewProductionEncoderConfig()
		level = zapcore.InfoLevel
	}
	// time layout 2006-01-02T15:04:05.000Z0700
	zcfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(zcfg)

	// colorized output
	zcfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(zcfg)

	// cores array
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	}

	if cfg.Dir != "" {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level))
	}

	// walk for tcp targets
	for _, target := range cfg.Targets {
		conn, err := net.Dial("tcp", target)
		if err != nil {
			log.Fatalf("error connecting to target(%s): %s", target, err)
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(conn), level))
	}

	core := zapcore.NewTee(cores...)

	// make new sugared logger
	// sugaredLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	sugaredLogger := zap.New(core).Sugar()
	if cfg.Debug {
		sugaredLogger = sugaredLogger.WithOptions(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	}
	sugaredLogger = sugaredLogger.Named(serviceName)
	if Logger != nil {
		sugaredLogger.Infof("logger initialized with targets: %s", cfg.Targets)
	}

	return &zapLogger{
		logger: sugaredLogger,
	}
}

func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *zapLogger) WithFields(args ...interface{}) logger.Logger {
	zl := zapLogger{
		logger: l.logger.With(args...),
	}
	return &zl
}
