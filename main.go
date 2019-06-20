package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

var essentialFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "port",
		Value: "8080",
		Usage: "Start HTTP server on port `PORT`",
	},
	cli.StringFlag{
		Name:  "grpc-port",
		Value: "8081",
		Usage: "Start grpc on port `GRPC_PORT`",
	},
	cli.IntFlag{
		Name:  "nats-port",
		Value: 4222,
		Usage: "Start nats on port `NATS_PORT`",
	},
	cli.IntFlag{
		Name:  "cluster-port",
		Value: 4248,
		Usage: "Start nats on port `NATS_PORT`",
	},
	cli.StringFlag{
		Name:  "config",
		Value: "none",
		Usage: "Load space cloud config from `FILE`",
	},
	cli.BoolFlag{
		Name:   "prod",
		Usage:  "Run space-cloud in production mode",
		EnvVar: "PROD",
	},
	cli.BoolFlag{
		Name:   "disable-metrics",
		Usage:  "Disable anonymous metric collection",
		EnvVar: "DISABLE_METRICS",
	},
	cli.BoolFlag{
		Name:   "disable-nats",
		Usage:  "Disable embedded nats server",
		EnvVar: "DISABLE_NATS",
	},
	cli.StringFlag{
		Name:   "seeds",
		Value:  "none",
		Usage:  "Seed nodes to cluster with",
		EnvVar: "SEEDS",
	},
}

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
			Flags:  essentialFlags,
		},
		{
			Name:   "start",
			Usage:  "runs the space cloud instance with mission control ui",
			Action: actionStart,
			Flags:  essentialFlags,
		},
		{
			Name:   "init",
			Usage:  "creates a config file with sensible defaults",
			Action: actionInit,
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
	grpcPort := c.String("grpc-port")
	natsPort := c.Int("nats-port")
	clusterPort := c.Int("cluster-port")
	configPath := c.String("config")
	isProd := c.Bool("prod")
	disableMetrics := c.Bool("disable-metrics")
	disableNats := c.Bool("disable-nats")
	seeds := c.String("seeds")

	// Project and env cannot be changed once space cloud has started
	s := initServer(isProd)

	if !disableNats {
		err := s.runNatsServer(seeds, natsPort, clusterPort)
		if err != nil {
			return err
		}
		fmt.Println("Started nats server on port ", defaultNatsOptions.Port)
	}

	if configPath != "none" {
		// Load the configFile from path if provided
		conf, err := config.LoadConfigFromFile(configPath)
		if err != nil {
			return err
		}

		// Save the config file path for future use
		s.configFilePath = configPath

		// Configure all modules
		err = s.loadConfig(conf)
		if err != nil {
			return err
		}
	}

	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.routineMetrics()
	}

	s.routes()
	return s.start(port, grpcPort)
}

func actionStart(c *cli.Context) error {
	// Load cli flags
	port := c.String("port")
	grpcPort := c.String("grpc-port")
	natsPort := c.Int("nats-port")
	clusterPort := c.Int("cluster-port")
	configPath := c.String("config")
	isProd := c.Bool("prod")
	disableMetrics := c.Bool("disable-metrics")
	disableNats := c.Bool("disable-nats")
	seeds := c.String("seeds")

	// Project and env cannot be changed once space cloud has started
	s := initServer(isProd)

	if !disableNats {
		err := s.runNatsServer(seeds, natsPort, clusterPort)
		if err != nil {
			return err
		}
		fmt.Println("Started nats server on port ", defaultNatsOptions.Port)
	}

	var conf *config.Project
	if configPath != "none" {
		// Load the configFile from path if provided
		config, err := config.LoadConfigFromFile(configPath)
		if err != nil {
			return err
		}

		conf = config
	} else {
		// Generate the config
		path, err := config.GenerateConfig(true, configPath)
		if err != nil {
			return nil
		}
		config, err := config.LoadConfigFromFile(path)
		if err != nil {
			return err
		}

		conf = config
		configPath = path
	}

	// Save the config file path for future use
	s.configFilePath = configPath

	// Configure all modules
	err := s.loadConfig(conf)
	if err != nil {
		return err
	}
	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.routineMetrics()
	}
	initMissionContol()
	s.routes()
	return s.start(port, grpcPort)
}

func actionInit(*cli.Context) error {
	_, err := config.GenerateConfig(false, "none")
	return err
}

func initMissionContol() error {
	homeDir := utils.UserHomeDir()
	uiPath := homeDir + "/space-cloud/mission-control-v" + buildVersion
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		resp, err := http.Get("https://spaceuptech.com/downloads/mission-control/mission-control-v" + buildVersion + ".zip")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(uiPath + ".zip")
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		err = utils.Unzip(uiPath+".zip", uiPath)
		if err != nil {
			return err
		}
		err = os.Remove(uiPath + ".zip")
		if err != nil {
			return err
		}
	}
	return nil
}
