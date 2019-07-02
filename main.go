package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/server"
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
		Value:  "",
		Usage:  "Seed nodes to cluster with",
		EnvVar: "SEEDS",
	},
	cli.BoolFlag{
		Name:   "profiler",
		Usage:  "Enable profiler endpoints for profiling",
		EnvVar: "PROFILER",
	},
}

func main() {
	app := cli.NewApp()
	app.Version = utils.BuildVersion
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
	profiler := c.Bool("profiler")

	// Project and env cannot be changed once space cloud has started
	s := server.New(isProd)

	if !disableNats {
		err := s.RunNatsServer(seeds, natsPort, clusterPort)
		if err != nil {
			return err
		}
	}

	if configPath != "none" {
		// Load the configFile from path if provided
		conf, err := config.LoadConfigFromFile(configPath)
		if err != nil {
			return err
		}

		// Save the config file path for future use
		s.SetConfigFilePath(configPath)

		// Configure all modules
		err = s.LoadConfig(conf)
		if err != nil {
			return err
		}
	}

	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.RoutineMetrics()
	}

	s.Routes(profiler)
	return s.Start(port, grpcPort, seeds)
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
	profiler := c.Bool("profiler")

	// Project and env cannot be changed once space cloud has started
	s := server.New(isProd)

	if configPath == "none" {
		configPath = "./config.yaml"
	}
	// Load config file
	conf, err := config.LoadConfigFromFile(configPath)
	// If config file does not exists then trigger the generate config flow
	if err != nil {
		err := config.GenerateConfig(configPath, true)
		if err != nil {
			return nil
		}
		conf, err = config.LoadConfigFromFile(configPath)
		if err != nil {
			return err
		}
	}

	// Save the config file path for future use
	s.SetConfigFilePath(configPath)

	// Configure all modules
	err = s.LoadConfig(conf)
	if err != nil {
		return err
	}

	// Anonymously collect usage metrics if not explicitly disabled
	if !disableMetrics {
		go s.RoutineMetrics()
	}

	err = initMissionContol()
	if err != nil {
		return err
	}

	fmt.Println("Started Mission Control UI at http://localhost:" + port + "/mission-control")

	if !disableNats {
		err := s.RunNatsServer(seeds, natsPort, clusterPort)
		if err != nil {
			return err
		}
		fmt.Println("Started nats server on port ", server.DefaultNatsOptions.Port)
	}

	s.Routes(profiler)
	return s.Start(port, grpcPort, seeds)
}

func actionInit(*cli.Context) error {
	return config.GenerateConfig("none", false)
}

func initMissionContol() error {
	homeDir := utils.UserHomeDir()
	uiPath := homeDir + "/.space-cloud/mission-control-v" + utils.BuildVersion
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		if _, err := os.Stat(homeDir + "/space-cloud"); os.IsNotExist(err) {
			os.Mkdir(homeDir+"/.space-cloud", os.ModePerm)
		}
		fmt.Println("Downloading Mission Control UI...")
		err := utils.DownloadFileFromURL("https://spaceuptech.com/downloads/mission-control/mission-control-v"+utils.BuildVersion+".zip", uiPath+".zip")
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
