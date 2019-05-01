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

func unitParseConnPQ (t *testing.T, cfg *config.ConnConfig,
		shouldFail bool) {
	fmt.Printf("testing ParseConnPostgres for %#v...\n", cfg)
	conn, err := ParseConnPostgres(cfg)
	if shouldFail && err == nil {
		t.Errorf("ParseConnPostgres should fail for %v", cfg)
	} else {
		fmt.Printf("conn: '%s'\n", conn)
	}
}//-- end func unitParseConnPQ

func TestParseConnPQ (t *testing.T) {
	fmt.Print("Testing ParseConnPostgres:\n")
	os.Setenv("TESTPASS", "test_password")
	unitParseConnPQ(t, &Config{
		User: "user",
		Auth: "ENV:TESTPASS",
		DBName: "tester_db",
		Protocol: "tcp",
		Host: "localhost", Port: "5432"}, false)
	fmt.Print("Done testing ParseConnPostgres\n\n")
}//-- end func TestParseConn

