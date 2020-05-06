package utils

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

// LogError logs the error in the proper format
func LogError(message, module, segment string, err error) error {
	// Prepare the fields
	entry := logrus.WithFields(logrus.Fields{"module": module, "segment": segment})
	if err != nil {
		entry = entry.WithError(err)
	}

	// Log the error
	entry.Errorln(message)

	// Return the error message
	return errors.New(message)
}

// LogWarn logs the warning message in the proper format
func LogWarn(message, module, segment string) {
	logrus.WithFields(logrus.Fields{"module": module, "segment": segment}).Warnln(message)
}

// LogInfo logs the info message in the proper format
func LogInfo(message, module, segment string) {
	logrus.WithFields(logrus.Fields{"module": module, "segment": segment}).Infoln(message)
}

// LogDebug logs the debug message in proper format
func LogDebug(message, module, segment string, extraFields map[string]interface{}) {
	entry := logrus.WithFields(logrus.Fields{"module": module, "segment": segment})
	if extraFields != nil {
		entry = entry.WithFields(extraFields)
	}
	entry.Debugln(message)
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
		_ = LogError(fmt.Sprintf("Invalid log level (%s) provided", loglevel), "utils", "set-log-level", nil)
		LogInfo("Defaulting to `info` level", "utils", "set-log-level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
