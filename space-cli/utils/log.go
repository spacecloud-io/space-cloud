package utils

import (
	"errors"

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
func LogDebug(message string) {
	logrus.Debugln(message)
}
