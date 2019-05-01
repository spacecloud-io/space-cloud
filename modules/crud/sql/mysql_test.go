package sql

/**
 * @author ollykel, ollykel416@gmail.com
 * @date Apr 22, 2019
 */

import (
	"fmt"
	"os"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
)

func unitParseConn (t *testing.T, cfg *config.ConnConfig, shouldFail bool) {
	fmt.Printf("testing ParseConnMySQL for %#v...\n", cfg)
	conn, err := ParseConnMySQL(cfg)
	if shouldFail && err == nil {
		t.Errorf("ParseConnMySQL should fail for %v", cfg)
	} else {
		fmt.Printf("conn: '%s'\n", conn)
	}
}//-- end func unitParseConn

func TestParseConn (t *testing.T) {
	fmt.Print("Testing ParseConnMySQL:\n")
	os.Setenv("TESTPASS", "test_password")
	unitParseConn(t, &Config{
		User: "user",
		Auth: "ENV:TESTPASS",
		DBName: "tester_db",
		Protocol: "tcp",
		Host: "localhost", Port: "3306"}, false)
	fmt.Print("Done testing ParseConnMySQL\n\n")
}//-- end func TestParseConn

