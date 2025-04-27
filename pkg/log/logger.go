package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger creates a new logger instance
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set log level from environment variable if available
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err == nil {
			logger.SetLevel(level)
		}
	} else {
		// Default to info level
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
