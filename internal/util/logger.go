package util

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(level string) *Logger {
	logger := logrus.New()

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Warnf("Invalid log level '%s': %v, defaulting to info level", level, err)
		parsedLevel = logrus.InfoLevel
	}

	logger.SetLevel(parsedLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		ForceColors:     true,
	})

	return &Logger{
		Logger: logger,
	}
}
