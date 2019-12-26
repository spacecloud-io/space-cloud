package main

import (
	"fmt"
	"log"
	"os"

	"github.com/segmentio/ksuid"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/metrics"
	"github.com/spaceuptech/space-cloud/utils/server"
)

var essentialFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "id",
		Value:  "none",
		Usage:  "The id to start space cloud with",
		EnvVar: "NODE_ID",
	},
	cli.StringFlag{
		Name:   "config",
		Value:  "config.yaml",
		Usage:  "Load space cloud config from `FILE`",
		EnvVar: "CONFIG",
	},
	cli.BoolFlag{
		Name:   "dev",
		Usage:  "Run space-cloud in development mode",
		EnvVar: "DEV",
	},
	cli.BoolFlag{
		Name:   "disable-metrics",
		Usage:  "Disable anonymous metric collection",
		EnvVar: "DISABLE_METRICS",
	},
	cli.BoolFlag{
		Name:   "profiler",
		Usage:  "Enable profiler endpoints for profiling",
		EnvVar: "PROFILER",
	},
	cli.StringFlag{
		Name:   "cluster",
		Usage:  "The cluster id to start space-cloud with",
		EnvVar: "CLUSTER_ID",
		Value:  "default-cluster",
	},
	cli.StringFlag{
		Name:   "advertise-addr",
		Usage:  "The address which will be broadcast to other space cloud instances",
		EnvVar: "ADVERTISE_ADDR",
	},
	cli.StringFlag{
		Name:   "store-type",
		Usage:  "The config store to use for storing project configs and other meta data",
		EnvVar: "STORE_TYPE",
		Value:  "none",
	},
	cli.IntFlag{
		Name:   "port",
		EnvVar: "PORT",
		Value:  4122,
	},
	cli.BoolFlag{
		Name:   "remove-project-scope",
		Usage:  "Removes the project level scope in the database and file storage modules",
		EnvVar: "REMOVE_PROJECT_SCOPE",
	},

	// Flags for the metrics module
	cli.BoolFlag{
		Name:   "enable-metrics",
		Usage:  "Enable the metrics module",
		EnvVar: "ENABLE_METRICS",
	},
	cli.BoolFlag{
		Name:   "disable-bandwidth",
		Usage:  "disable the bandwidth measurement",
		EnvVar: "DISABLE_BANDWIDTH",
	},
	cli.StringFlag{
		Name:   "metrics-sink",
		Usage:  "The sink to output metrics data to",
		EnvVar: "METRICS_SINK",
	},
	cli.StringFlag{
		Name:   "metrics-conn",
		Usage:  "The connection string of the sink",
		EnvVar: "METRICS_CONN",
	},
	cli.StringFlag{
		Name:   "metrics-scope",
		Usage:  "The database / topic to push the metrics to",
		EnvVar: "METRICS_SCOPE",
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
	nodeID := c.String("id")
	configPath := c.String("config")
	isDev := c.Bool("dev")
	disableMetrics := c.Bool("disable-metrics")
	disableBandwidth := c.Bool("disable-bandwidth")
	profiler := c.Bool("profiler")

	// Load flag related to the port
	port := c.Int("port")
	removeProjectScope := c.Bool("remove-project-scope")

	// Load flags related to ssl

	// Flags related to the admin details

	// Load flags related to clustering
	clusterID := c.String("cluster")
	storeType := c.String("store-type")
	advertiseAddr := c.String("advice-addr")

	// Load the flags for the metrics module
	enableMetrics := c.Bool("enable-metrics")
	metricsSink := c.String("metrics-sink")
	metricsConn := c.String("metrics-conn")
	metricsScope := c.String("metrics-scope")

	// Generate a new id if not provided
	if nodeID == "none" {
		nodeID = "auto-" + ksuid.New().String()
	}

	s, err := server.New(nodeID, clusterID, advertiseAddr, storeType, removeProjectScope,
		&metrics.Config{IsEnabled: enableMetrics, SinkType: metricsSink, SinkConn: metricsConn, Scope: metricsScope, DisableBandwidth: disableBandwidth})
	if err != nil {
		return err
	}

	// Load the configFile from path if provided
	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		conf = config.GenerateEmptyConfig()
	}

	// Save the config file path for future use
	s.SetConfigFilePath(configPath)

	// Download and host mission control
	staticPath, err := initMissionContol(utils.BuildVersion)
	if err != nil {
		return err
	}

	// Initialise the routes
	s.InitRoutes(profiler, staticPath)

	// Configure all modules
	s.SetConfig(conf, !isDev)

	return s.Start(disableMetrics, port)
}

func actionInit(*cli.Context) error {
	return config.GenerateConfig("none")
}

func initMissionContol(version string) (string, error) {
	homeDir := utils.UserHomeDir()
	uiPath := homeDir + "/.space-cloud/mission-control-v" + version
	_, err := os.Stat(uiPath)
	if os.IsNotExist(err) {
		fmt.Println("Could not find mission control")
		_, err := os.Stat(homeDir + "/.space-cloud")
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		if os.IsNotExist(err) {
			err := os.Mkdir(homeDir+"/.space-cloud", os.ModePerm)
			if err != nil {
				return "", err
			}
		}
		fmt.Println("Downloading...")
		err = utils.DownloadFileFromURL("https://spaceuptech.com/downloads/mission-control/mission-control-v"+version+".zip", uiPath+".zip")
		if err != nil {
			return "", err
		}
		fmt.Println("Extracting...")
		err = utils.Unzip(uiPath+".zip", uiPath)
		if err != nil {
			return "", err
		}
		fmt.Println("Done...")
		err = os.Remove(uiPath + ".zip")
		if err != nil {
			return "", err
		}
		return uiPath + "/build", nil
	}
	if err != nil {
		return "", err
	}
	return uiPath + "/build", nil
}
