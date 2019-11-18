package stdOut

import (
	"fmt"
	"github.com/spaceuptech/space-cloud/utils"
)

type StdOut struct {
	enabled bool
	levels []utils.LogLevel
}

func Init(levels []utils.LogLevel, enabled bool) (*StdOut, error){
	stdOutRef := &StdOut{
		enabled: enabled,
		levels:  levels,
	}
	return stdOutRef, nil
}


func (stdOutRef StdOut) Debug(msg string, data map[string]interface{}) error {
	fmt.Println(utils.LogLevel(utils.DEBUG).String(), ":", msg, data)
	return nil
}
func (stdOutRef StdOut) Info(msg string, data map[string]interface{}) error{
	fmt.Println(utils.LogLevel(utils.INFO).String(), ":", msg, data)
	return nil
}
func (stdOutRef StdOut) Warning(msg string, data map[string]interface{}) error{
	fmt.Println(utils.LogLevel(utils.WARNING).String(), ":", msg, data)
	return nil
}
func (stdOutRef StdOut) Error(msg string, data map[string]interface{}) error{
	fmt.Println(utils.LogLevel(utils.ERROR).String(), ":", msg, data)
	return nil
}

func (stdOutRef StdOut) Close () error{
	return nil
}