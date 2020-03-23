package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/server"
)

const (
	loglevelDebug = "debug"
	loglevelInfo  = "info"
	logLevelError = "error"
)

var essentialFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "log-level",
		EnvVar: "LOG_LEVEL",
		Usage:  "Set the log level [debug | info | error]",
		Value:  loglevelInfo,
	},
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
		Value:  "localhost:4122",
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
	cli.StringFlag{
		Name:   "restrict-hosts",
		EnvVar: "RESTRICT_HOSTS",
		Usage:  "Comma separated values of the hosts to restrict mission-control to",
		Value:  "*",
	},
	cli.BoolFlag{
		Name:   "remove-project-scope",
		Usage:  "Removes the project level scope in the database and file storage modules",
		EnvVar: "REMOVE_PROJECT_SCOPE",
	},
	cli.StringFlag{
		Name:   "runner-addr",
		Usage:  "The address used to reach the runner",
		EnvVar: "RUNNER_ADDR",
	},
	cli.StringFlag{
		Name:   "artifact-addr",
		Usage:  "The address used to reach the artifact server",
		EnvVar: "ARTIFACT_ADDR",
	},

	// Flags for ssl
	cli.BoolFlag{
		Name:   "ssl-enable",
		Usage:  "Enable https and lets encrypt support",
		EnvVar: "SSL_ENABLE",
	},
	cli.StringFlag{
		Name:   "ssl-cert",
		Value:  "none",
		Usage:  "Load ssl certificate from `FILE`",
		EnvVar: "SSL_CERT",
	},
	cli.StringFlag{
		Name:   "ssl-key",
		Value:  "none",
		Usage:  "Load ssl key from `FILE`",
		EnvVar: "SSL_KEY",
	},

	// flags for admin man
	cli.StringFlag{
		Name:   "admin-user",
		Usage:  "Set the admin user name",
		EnvVar: "ADMIN_USER",
		Value:  "",
	},
	cli.StringFlag{
		Name:   "admin-pass",
		Usage:  "Set the admin password",
		EnvVar: "ADMIN_PASS",
		Value:  "",
	},
	cli.StringFlag{
		Name:   "admin-secret",
		Usage:  "Set the admin secret",
		EnvVar: "ADMIN_SECRET",
		Value:  "",
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
	logLevel := c.String("log-level")
	setLogLevel(logLevel)

	// Load flag related to the port
	port := c.Int("port")

	removeProjectScope := c.Bool("remove-project-scope")
	runnerAddr := c.String("runner-addr")
	artifactAddr := c.String("artifact-addr")

	// Load flags related to ssl
	sslEnable := c.Bool("ssl-enable")
	sslKey := c.String("ssl-key")
	sslCert := c.String("ssl-cert")

	// Flags related to the admin details
	adminUser := c.String("admin-user")
	adminPass := c.String("admin-pass")
	adminSecret := c.String("admin-secret")

	// Load flags related to clustering
	clusterID := c.String("cluster")
	storeType := c.String("store-type")
	advertiseAddr := c.String("advertise-addr")

	// Load the flags for the metrics module
	enableMetrics := c.Bool("enable-metrics")
	metricsSink := c.String("metrics-sink")
	metricsConn := c.String("metrics-conn")
	metricsScope := c.String("metrics-scope")

	// Generate a new id if not provided
	if nodeID == "none" {
		nodeID = "auto-" + ksuid.New().String()
	}

	s, err := server.New(nodeID, clusterID, advertiseAddr, storeType, runnerAddr, artifactAddr, removeProjectScope,
		&metrics.Config{IsEnabled: enableMetrics, SinkType: metricsSink, SinkConn: metricsConn, Scope: metricsScope, DisableBandwidth: disableBandwidth})
	if err != nil {
		return err
	}

	// Load the configFile from path if provided
	conf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		conf = config.GenerateEmptyConfig()
	}
	if conf.Admin == nil {
		conf.Admin = config.GenerateAdmin()
	}
	// Save the config file path for future use
	s.SetConfigFilePath(configPath)

	// Override the admin config if provided
	if adminUser != "" {
		conf.Admin.Users[0].User = adminUser
	}
	if adminPass != "" {
		conf.Admin.Users[0].Pass = adminPass
	}
	if adminSecret != "" {
		conf.Admin.Secret = adminSecret
	}

	// Download and host mission control
	staticPath, err := initMissionContol(utils.BuildVersion)
	if err != nil {
		return err
	}

	// Set the ssl config
	if sslEnable {
		conf.SSL = &config.SSL{Enabled: true, Crt: sslCert, Key: sslKey}
	}

	// Configure all modules
	s.SetConfig(conf, !isDev)

	return s.Start(profiler, disableMetrics, staticPath, port, strings.Split(c.String("restrict-hosts"), ","))
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

func setLogLevel(loglevel string) {
	switch loglevel {
	case loglevelDebug:
		logrus.SetLevel(logrus.DebugLevel)
	case loglevelInfo:
		logrus.SetLevel(logrus.InfoLevel)
	case logLevelError:
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Errorf("Invalid log level (%s) provided", loglevel)
		logrus.Infoln("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
