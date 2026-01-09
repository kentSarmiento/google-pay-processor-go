package logger

import (
	"os"
	"time"

	"github.com/kentSarmiento/google-pay-processor-go/internal/domain"
	"github.com/rs/zerolog"
)

// ZerologLogger implements the domain.Logger interface using zerolog
type ZerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger creates a new zerolog logger
func NewZerologLogger(level, format string) domain.Logger {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	// Configure output format
	var logger zerolog.Logger
	if format == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(logLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		logger = zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return &ZerologLogger{logger: logger}
}

func (l *ZerologLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *ZerologLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *ZerologLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *ZerologLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.Error().Fields(fields).Msg(msg)
}

func (l *ZerologLogger) Fatal(msg string, fields map[string]interface{}) {
	l.logger.Fatal().Fields(fields).Msg(msg)
}
