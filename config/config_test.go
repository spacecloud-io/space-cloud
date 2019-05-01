package config

/**
 * @author ollykel, ollykel416@gmail.com
 * @date Apr 19, 2019
 * Testing functions for config package
 */

import (
	"testing"
	"fmt"
)

func unitTestGetSuffix (t *testing.T, fname, correct string) {
	fmt.Printf("testing filename: '%s', suffix: '%s'...\n", fname, correct)
	suffix := getSuffix(fname)
	if suffix != correct {
		t.Errorf("wanted %s, got %s", correct, suffix)
	}
}//-- end func unitTestGetSuffix

func TestGetSuffix (t *testing.T) {
	fmt.Print("Testing func getSuffix:\n")
	unitTestGetSuffix(t, "config.yaml", "yaml")
	unitTestGetSuffix(t, "./config.json", "json")
	unitTestGetSuffix(t, "../config.yml", "yml")
	unitTestGetSuffix(t, "github.com/main/config.xml", "xml")
	unitTestGetSuffix(t, "github.com/main/config", "")
	fmt.Print("Done testing func getSuffix\n\n")
}//-- end func TestGetSuffix

// LoadConfigFromFile

func unitLoadConfigFromFile (t *testing.T, fname string) {
	fmt.Printf("loading config from '%s'...\n", fname)
	proj, err := LoadConfigFromFile(fname)
	if err != nil { t.Errorf("error: %s", err.Error()) }
	fmt.Printf("Config:\n%#v\n", proj)
}//-- end func unitLoadConfigFromFile

func TestLoadConfigFromFile (t *testing.T) {
	fmt.Print("Testing LoadConfigFromFile:\n")
	unitLoadConfigFromFile(t, "./.test/test-conf.yaml")
	unitLoadConfigFromFile(t, "./.test/tester.yaml")
	fmt.Print("Done testing LoadConfigFromFile\n\n")
}//-- end func TestLoadConfigFromFile

