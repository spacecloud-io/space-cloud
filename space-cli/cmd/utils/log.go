package utils

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

// LogError logs the error in the proper format
func LogError(message string, err error) error {

	// Log with error if provided
	if err != nil {
		logrus.WithField("error", err.Error()).Errorln(message)
	} else {
		logrus.Errorln(message)
	}

	// Return the error message
	return errors.New(message)
}

// LogInfo logs te info message in the proper format
func LogInfo(message string) {
	logrus.Infoln(message)
}

// LogDebug logs the debug message in proper format
func LogDebug(message string, extraFields map[string]interface{}) {
	if extraFields != nil {
		logrus.WithFields(extraFields).Debugln(message)
		return
	}
	logrus.Debugln(message)
}

// SetLogLevel sets a single verbosity level for log messages.
func SetLogLevel(loglevel string) {
	switch loglevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		_ = LogError(fmt.Sprintf("Invalid log level (%s) provided", loglevel), nil)
		LogInfo("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
