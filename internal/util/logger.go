package util

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

func NewLogger(level, format string) *Logger {
	base := logrus.New()

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		base.Warnf("Invalid log level '%s': %v, defaulting to info", level, err)
		parsedLevel = logrus.InfoLevel
	}
	base.SetLevel(parsedLevel)

	if format == "json" {
		base.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	} else {
		base.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			ForceColors:     true,
		})
	}

	return &Logger{Entry: logrus.NewEntry(base)}
}

func (l *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{Entry: l.Entry.WithFields(fields)}
}
