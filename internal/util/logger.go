package util

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(level string) *logrus.Logger {
	logger := logrus.New()

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Warnf("Invalid log level '%s': %v, defaulting to info level", level, err)
		parsedLevel = logrus.InfoLevel
	}

	logger.SetLevel(parsedLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	return logger
}
