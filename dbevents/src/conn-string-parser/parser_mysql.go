package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func parseMySQLConn(conn string) error {
	c, err := mysql.ParseDSN(conn)
	if err != nil {
		return err
	}

	arr := strings.Split(c.Addr, ":")
	jsonString, _ := json.Marshal(DBConfig{
		Host:    arr[0],
		Port:    arr[1],
		User:    c.User,
		Pass:    c.Passwd,
		SSLMode: getMySQLSSLMode(c),
	})

	// Print it out
	fmt.Println(string(jsonString))
	return nil
}

func getMySQLSSLMode(c *mysql.Config) string {
	switch c.TLSConfig {
	case "true":
		return "verify-ca"
	case "false":
		return "disabled"
	case "skip-verify":
		return "required"
	case "preferred":
		return "preferred"
	}

	return "disabled"
}
