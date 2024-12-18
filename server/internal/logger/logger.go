package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitializeLogger sets up the logger
func InitializeLogger(env string) {
	Log = logrus.New()

	Log.SetOutput(os.Stdout)

	// Set log format
	if env == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{}) // JSON for production
	} else {
		Log.SetFormatter(&logrus.TextFormatter{ // Pretty for development
			FullTimestamp: true,
		})
	}

	if env == "production" {
		Log.SetLevel(logrus.InfoLevel) // Info and above for production
	} else {
		Log.SetLevel(logrus.DebugLevel) // Debug and above for development
	}
}

func GetLogger() *logrus.Logger {
	if Log == nil {
		InitializeLogger("development") // Default to "development" if not initialized
	}
	return Log
}
