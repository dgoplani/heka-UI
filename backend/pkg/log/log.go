package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

var formatterMap = map[string]logrus.Formatter{
	"json": &logrus.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05.999",
	},
	"text": &logrus.TextFormatter{},
}

// Setup - initialize logging functionality to be compatible with log collector
func Setup(format, level string) {
	logger = logrus.StandardLogger()
	if formatter, ok := formatterMap[format]; ok {
		logger.SetFormatter(formatter)
	} else {
		logger.Fatalf("Invalid log format: %v", format)
	}
	loglevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Fatalf("Invalid log level: %v", err)
	}
	logger.SetLevel(loglevel)
}

