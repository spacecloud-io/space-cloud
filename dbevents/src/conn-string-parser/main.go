package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "conn-string-parser",
		Version: "0.1.0",
		Commands: []*cli.Command{
			{
				Name:  "parse",
				Usage: "Parse the db connection string",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "db-type",
						Usage: "The `db` the connection string is for",
						Value: "none",
					},
				},
				Action: func(c *cli.Context) error {
					dbType := c.String("db-type")
					if dbType == "none" {
						return errors.New("db-type is a required flag")
					}
					if c.Args().Len() == 0 {
						return errors.New("connection string not provided")
					}
					return parseConnectionString(dbType, c.Args().First())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		jsonString, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Println(string(jsonString))
	}
}
