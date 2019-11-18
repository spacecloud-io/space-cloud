package logging

import (
	"errors"
	"fmt"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/logging/stdOut"
	"sync"
)

type Module struct {
	sync.RWMutex
	loggers map[string]Logger
}

type Logger interface {
	Debug(msg string, data map[string]interface{}) error
	Info(msg string, data map[string]interface{}) error
	Warning(msg string, data map[string]interface{}) error
	Error(msg string, data map[string]interface{}) error
	Close() error
}

func Init() *Module {
	return &Module{loggers: make(map[string]Logger)}
}

func initLogger(logType utils.LoggerType, levels []LogLevel, enabled bool) (Logger, error){
	switch logType {
	case utils.StdOut:
		return stdOut.Init(levels, enabled)
	}

	return nil, errors.New("given loggerType could not be resolved")
}

func (m *Module) SetConfig(log config.Log) error {
	m.Lock()
	defer m.Unlock()

	// Close the previous database connections
	for _, v := range m.loggers {
		v.Close()
	}
	m.loggers = make(map[string]Logger, len(log))

	// Create a new crud blocks
	for k, v := range log {
		c, err := initLogger(utils.LoggerType(k), v.Levels, v.Enabled)
		m.loggers[k] = c

		if err != nil {
			fmt.Println("Error creating Logger with key", k)
		}
	}
	return nil
}

func (m *Module) Debug(msg string, data map[string]interface{}) error {
	for _, v := range m.loggers{
		err :=v.Debug(msg, data); if err != nil {
			return err
		}
	}

	return nil
}
func (m Module) Info(msg string, data map[string]interface{}) error{
	for _, v := range m.loggers{
		err :=v.Info(msg, data); if err != nil {
			return err
		}
	}

	return nil
}
func (m Module) Warning(msg string, data map[string]interface{}) error{
	for _, v := range m.loggers{
		err :=v.Warning(msg, data); if err != nil {
			return err
		}
	}

	return nil
}
func (m Module) Error(msg string, data map[string]interface{}) error{
	for _, v := range m.loggers{
		err :=v.Error(msg, data); if err != nil {
			return err
		}
	}

	return nil
}

func (m Module) Close() error{
	for _, v := range m.loggers{
		err :=v.Close(); if err != nil {
			return err
		}
	}

	return nil
}