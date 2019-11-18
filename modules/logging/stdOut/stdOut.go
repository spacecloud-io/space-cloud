package stdOut

import (
	"fmt"
	"github.com/spaceuptech/space-cloud/utils/logging"
)

type StdOut struct {
	enabled bool
	levels []logging.LogLevel
}

func Init(levels []logging.LogLevel, enabled bool) (*StdOut, error){
	stdOutRef := &StdOut{
		enabled: enabled,
		levels:  levels,
	}
	return stdOutRef, nil
}


func (stdOutRef StdOut) Debug(msg string, data map[string]interface{}) error {
	fmt.Println(logging.LogLevel(logging.DEBUG).String(), ":", msg, data)
	return nil
}
func (stdOutRef StdOut) Info(msg string, data map[string]interface{}) error{
	fmt.Println(logging.LogLevel(logging.INFO).String(), ":", msg, data)
	return nil
}
func (stdOutRef StdOut) Warning(msg string, data map[string]interface{}) error{
	fmt.Println(logging.LogLevel(logging.WARNING).String(), ":", msg, data)
	return nil
}
func (stdOutRef StdOut) Error(msg string, data map[string]interface{}) error{
	fmt.Println(logging.LogLevel(logging.ERROR).String(), ":", msg, data)
	return nil
}

func (stdOutRef StdOut) Close () error{
	return nil
}