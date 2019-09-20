package main

import (
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
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
		Name:  "config",
		Value: "config.yaml",
		Usage: "Load space cloud config from `FILE`",
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
		Name:   "disable-nats",
		Usage:  "Disable embedded nats server",
		EnvVar: "DISABLE_NATS",
	},
	cli.StringFlag{
		Name:   "seeds",
		Value:  "127.0.0.1",
		Usage:  "Seed nodes to cluster with",
		EnvVar: "SEEDS",
	},
	cli.BoolFlag{
		Name:   "profiler",
		Usage:  "Enable profiler endpoints for profiling",
		EnvVar: "PROFILER",
	},
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
	disableNats := c.Bool("disable-nats")
	seeds := c.String("seeds")
	profiler := c.Bool("profiler")

	// Load flags related to ssl
	sslCert := c.String("ssl-cert")
	sslKey := c.String("ssl-key")

	// Flags related to the admin details
	adminUser := c.String("admin-user")
	adminPass := c.String("admin-pass")
	adminSecret := c.String("admin-secret")

	// Generate a new id if not provided
	if nodeID == "none" {
		nodeID = uuid.NewV1().String()
	}

	// Start nats if not disabled
	if !disableNats {
		err := server.RunNatsServer(seeds, utils.PortNatsServer, utils.PortNatsCluster)
		if err != nil {
			return err
		}
	}

	s, err := server.New(nodeID)
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

	// Initialise the routes
	s.InitRoutes(profiler, staticPath)

	// Set the ssl config
	if sslCert != "none" && sslKey != "none" {
		s.InitSecureRoutes(profiler, staticPath)
		conf.SSL = &config.SSL{Enabled: true, Crt: sslCert, Key: sslKey}
	}

	// Configure all modules
	s.SetConfig(conf, !isDev)

	return s.Start(seeds, disableMetrics)
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
