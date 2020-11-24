package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
)

func parseSQLSeverConnString(conn string) error {
	u, err := url.Parse(conn)
	if err != nil {
		return err
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return err
	}

	pass, _ := u.User.Password()
	jsonString, _ := json.Marshal(DBConfig{
		Host: host,
		Port: port,
		User: u.User.Username(),
		Pass: pass,
		DB:   u.Query().Get("database"),
	})

	// Print it out
	fmt.Println(string(jsonString))
	return nil
}
