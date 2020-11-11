package main

import (
	"encoding/json"
	"errors"
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

	// Throw error if path is empty
	if u.Path == "" {
		return errors.New("no database provided in conn string")
	}

	pass, _ := u.User.Password()
	jsonString, _ := json.Marshal(DBConfig{
		Host: host,
		Port: port,
		User: u.User.Username(),
		Pass: pass,
		DB:   u.Path[1:],
	})

	// Print it out
	fmt.Println(string(jsonString))
	return nil
}
