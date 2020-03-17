package utils

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// LogError logs the error in the proper format
func LogError(message, module, segment string, err error) error {
	entry := logrus.WithField("module", module)
	entry = entry.WithField("segment", segment)
	if err != nil {
		entry = entry.WithField("error", err.Error())
	}
	entry.Errorln(message)
	return errors.New(message)
}

// LogDebug logs the debug message in proper format
func LogDebug(message, module, segment string, extraFields map[string]interface{}) {
	entry := logrus.WithFields(logrus.Fields{"module": module, "segment": segment})
	if extraFields != nil {
		entry = entry.WithFields(extraFields)
	}
	entry.Debugln(message)
}
