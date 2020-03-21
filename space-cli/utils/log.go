package utils

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// LogError logs the error in the proper format
func LogError(message, module, segment string, err error) error {
	entry := logrus.WithField("module", module)

	// Add segment if provided
	if segment != "" {
		entry = entry.WithField("segment", segment)
	}

	// Add error if provided
	if err != nil {
		entry = entry.WithField("error", err.Error())
	}

	// Log the message
	entry.Errorln(message)

	// Return the error message
	return errors.New(message)
}

// LogInfo logs te info message in the proper format
func LogInfo(message, module, segment string) {
	logrus.WithFields(map[string]interface{}{"module": module, "segment": segment}).Infoln(message)
}

// LogDebug logs the debug message in proper format
func LogDebug(message, module, segment string, extraFields map[string]interface{}) {
	entry := logrus.WithFields(logrus.Fields{"module": module, "segment": segment})
	if extraFields != nil {
		entry = entry.WithFields(extraFields)
	}
	entry.Debugln(message)
}
