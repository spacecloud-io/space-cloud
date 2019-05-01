package config

import (
	"encoding/json"
	"os"
	"io"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type decoder interface {
	Decode (dest interface{}) error
}//-- end decoder interface

func getSuffix (s string) string {
	lastPathIdx := strings.LastIndex(s, string(os.PathSeparator))
	var path string
	if lastPathIdx < 0 {
		path = s
	} else {
		path = s[lastPathIdx + 1:]
	}
	lastIndex := strings.LastIndexByte(path, '.')
	if lastIndex < 0 || lastIndex > len(path) - 2 { return "" }
	return path[lastIndex + 1:]
}//-- end func getSuffix

func LoadConfig (src io.Reader, dataType string) (*Project, error) {
	var dec decoder
	switch dataType {
		case "yaml", "yml":
			dec = yaml.NewDecoder(src)
		case "json":
			dec = json.NewDecoder(src)
		default:
			return nil, fmt.Errorf("unrecognized file type %s", dataType)
	}//-- end switch
	proj := new(Project)
	if err := dec.Decode(proj); err != nil { return nil, err }
	return proj, nil
}//-- end func LoadConfig

// LoadConfigFromFile loads the config from the provided file path
func LoadConfigFromFile(path string) (*Project, error) {
	file, err := os.Open(path); if err != nil { return nil, err }
	defer file.Close()

	return LoadConfig(file, getSuffix(path))
}//-- end func LoadConfigFromFile

func (proj *Project) WriteTo (dest io.Writer) error {
	encoder := yaml.NewEncoder(dest); defer encoder.Close()
	return  encoder.Encode(proj)
}//-- end func Project.Write

func (proj *Project) Save () error {
	fname := "./" + proj.ID + ".yaml"
	f, err := os.Create(fname); if err != nil { return err }
	defer f.Close()
	f.WriteString("---\n")
	err = proj.WriteTo(f); if err != nil { return err }
	return f.Sync()
}//-- end func Project.Save

