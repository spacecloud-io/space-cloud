package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/config"
)

func main() {
	app := cli.NewApp()
	app.Version = buildVersion
	app.Name = "space-cloud"
	app.Usage = "core binary to run space cloud"

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "runs the space cloud instance",
			Action: actionRun,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "8080",
					Usage: "Start HTTP server on port `PORT`",
				},
				cli.StringFlag{
					Name:  "config",
					Value: "none",
					Usage: "Load space cloud config from `FILE`",
				},
				cli.BoolFlag{
					Name:  "prod",
					Usage: "The environment to start space cloud in",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func actionRun(c *cli.Context) error {
	// Load cli flags
	port := c.String("port")
	configPath := c.String("config")
	isProd := c.Bool("env")

	// Project and env cannot be changed once space cloud has started
	s := initServer(isProd)

	if configPath != "none" {
		config, err := config.LoadConfigFromFile(configPath)
		if err != nil {
			return err
		}
		err = s.loadConfig(config)
		if err != nil {
			return err
		}
	}

	s.routes()
	return s.start(port)
}
