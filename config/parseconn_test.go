package config

/**
 * @author ollykel, ollykel416@gmail.com
 * @date Apr 22, 2019
 */

import (
	"fmt"
	"os"
	"testing"
)

func unitParseConn (t *testing.T, cfg *ConnConfig, shouldFail bool) {
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
	unitParseConn(t, &ConnConfig{
		User: "user",
		Auth: "ENV:TESTPASS",
		DBName: "tester_db",
		Protocol: "tcp",
		Host: "localhost", Port: "3306"}, false)
	fmt.Print("Done testing ParseConnMySQL\n\n")
}//-- end func TestParseConn

func unitParseConnPQ (t *testing.T, cfg *ConnConfig, shouldFail bool) {
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
	unitParseConnPQ(t, &ConnConfig{
		User: "user",
		Auth: "ENV:TESTPASS",
		DBName: "tester_db",
		Protocol: "tcp",
		Host: "localhost", Port: "5432"}, false)
	fmt.Print("Done testing ParseConnPostgres\n\n")
}//-- end func TestParseConn

func unitParseConnMgo (t *testing.T, cfg *ConnConfig, shouldFail bool) {
	fmt.Printf("testing ParseConnMongo for %#v...\n", cfg)
	conn, err := ParseConnMongo(cfg)
	if shouldFail && err == nil {
		t.Errorf("ParseConnMongo should fail for %v", cfg)
	} else {
		fmt.Printf("conn: '%s'\n", conn)
	}
}//-- end func unitParseConnMgo

func TestParseConnMgo (t *testing.T) {
	fmt.Print("Testing ParseConnMongo:\n")
	os.Setenv("TESTPASS", "weewoo")
	unitParseConnMgo(t, &ConnConfig{
		User: "mongo_user",
		Auth: "ENV:TESTPASS",
		DBName: "tester_db",
		Protocol: "tcp",
		Host: "mongo.example.com", Port: "27017"}, false)
	fmt.Print("Done testing ParseConnPostgres\n\n")
}//-- end func TestParseConn

