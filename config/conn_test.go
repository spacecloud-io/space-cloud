package config

/**
 * @author ollykel416, ollykel416@gmail.com
 * @date Apr 22, 2019
 * Testers for ConnConfig generating functions in sql/config.go
 */

import (
	"fmt"
	"path/filepath"
	"testing"
	"os"
)

const testDir = ".test"

// filenames should be under local .test dir
func unitGetPassword (t *testing.T, auth, correct string,
		shouldFail bool) {
	fmt.Printf("testing GetPassword for %s...\n", auth)
	pass, err := GetPassword(auth)
	if !shouldFail && err != nil {
		t.Errorf("shouldn't fail: '%v'", err)
	} else if shouldFail && err == nil {
		t.Errorf("should fail for '%s'", auth)
	}
	if pass != correct && !shouldFail {
		t.Errorf("wanted: '%s', got: '%s'", correct, pass)
	}
}//-- end func unitReadPassword

func makePasswordFileAuth (fname string) string {
	return "FILE:" + filepath.Join(testDir, fname)
}//-- end func makePasswordFileAuth

func TestGetPassword (t *testing.T) {
	fmt.Print("Testing func GetPassword:\n")
	unitGetPassword(t, makePasswordFileAuth("pass-1.txt"), "weewoo", false)
	// non-existant file
	unitGetPassword(t, makePasswordFileAuth("nothing.txt"), "", true)
	testPass := "p@ssword101"
	os.Setenv("TEST_PASS", testPass)
	unitGetPassword(t, "ENV:TEST_PASS", testPass, false)
	// non-existant env var
	unitGetPassword(t, "ENV:NOTHING", "", true)
	// strings
	unitGetPassword(t, "STRING:p@ssword", "p@ssword", false)
	fmt.Print("Done testing func GetPassword\n\n")
}//-- end func TestGetPassword

