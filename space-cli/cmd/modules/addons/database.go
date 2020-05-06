package addons

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cli/cmd/modules/operations"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

type loadEnvResponse struct {
	Quotas struct {
		MaxDatabases int `json:"maxDatabases"`
		MaxProjects  int `json:"maxProjects"`
	} `json:"quotas"`
}

func addDatabase(dbtype, username, password, alias, version string) error {
	ctx := context.Background()
	autoApply := viper.GetBool("auto-apply")
	project := viper.GetString("project")
	if autoApply {
		if project == "" {
			return utils.LogError(`Please provide project id through "--project" flag`, nil)
		}
		// fetch quotas from gateway
		resp := new(loadEnvResponse)
		err := utils.Get(http.MethodGet, "/v1/config/env", map[string]string{}, resp)
		if err != nil {
			return utils.LogError(`Cannot fetch quotas from gateway, Is gateway running ?`, err)
		}

		// fetch current db config
		dbConfig, err := database.GetDbConfig(project, "db-config", map[string]string{})
		if err != nil {
			return utils.LogError(`Cannot fetch database config from gateway`, err)
		}

		// check if database can be added
		if (len(dbConfig) + 1) > resp.Quotas.MaxDatabases {
			return utils.LogError(fmt.Sprintf(`Cannot add database in project "%s", max database limit reached. upgrade you plan`, project), err)
		}
	}

	// Prepare the docker image name name
	dockerImage := fmt.Sprintf("%s:%s", dbtype, version)

	// Set alias if not provided
	if alias == "" {
		alias = dbtype
	}

	// Set the environment variables
	var env []string
	switch dbtype {
	case "mysql":
		if username != "root" {
			return utils.LogError("Only the username root is allowed for MySQL", nil)
		}
		env = []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password)}
	case "postgres":
		env = []string{fmt.Sprintf("POSTGRES_USER=%s", username), fmt.Sprintf("POSTGRES_PASSWORD=%s", password)}
	case "mongo":
		if username != "" || password != "" {
			return utils.LogError("Cannot set username or password with Mongo", nil)
		}
		env = []string{}
	default:
		return utils.LogError(fmt.Sprintf("Invalid database type (%s) provided", dbtype), nil)
	}

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Check if a database container already exist
	filterArgs := filters.Arg("label", "app=space-cloud")
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if database already exists", err)
	}
	if len(containers) == 0 {
		utils.LogInfo("No space-cloud instance found. Run 'space-cli setup' first")
		return nil
	}

	// Pull image if it doesn't already exist
	if err := utils.PullImageIfNotExist(ctx, docker, dockerImage); err != nil {
		return utils.LogError(fmt.Sprintf("Could not pull the image (%s). Make sure docker is running and that you have an active internet connection.", dockerImage), err)
	}

	// Create the database
	containerRes, err := docker.ContainerCreate(ctx, &container.Config{
		Labels: map[string]string{"app": "addon", "service": "db", "name": alias},
		Image:  dockerImage,
		Env:    env,
	}, &container.HostConfig{
		NetworkMode: "space-cloud",
	}, nil, fmt.Sprintf("space-cloud--addon--db--%s", alias))
	if err != nil {
		return utils.LogError("Unable to create local docker database", err)
	}

	// Start the database
	if err := docker.ContainerStart(ctx, containerRes.ID, types.ContainerStartOptions{}); err != nil {
		return utils.LogError("Unable to start local docker database", err)
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to open hosts file", err)
	}

	// Get the container's info
	info, err := docker.ContainerInspect(ctx, containerRes.ID)
	if err != nil {
		return utils.LogError(fmt.Sprintf("Unable to inspect c (%s)", containerRes.ID), err)
	}

	hostName := utils.GetServiceDomain("db", alias)

	// Remove the domain from the hosts file
	hosts.RemoveHost(hostName)

	// Add it back with the new ip address
	hosts.AddHost(info.NetworkSettings.Networks["space-cloud"].IPAddress, hostName)

	// Save the hosts file
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", err)
	}

	if autoApply {
		connDefault := ""
		duration := 15
		logrus.Println("host", hostName)
		switch dbtype {
		case "postgres":
			connDefault = fmt.Sprintf("postgres://postgres:mysecretpassword@%s:5432/postgres?sslmode=disable", hostName)
		case "mongo":
			connDefault = fmt.Sprintf("mongodb://%s:27017", hostName)
		case "mysql":
			connDefault = fmt.Sprintf("root:my-secret-pw@tcp(%s:3306)/", hostName)
			duration = 220
		default:
			return fmt.Errorf("invalid database provided, supported databases postgres,sqlserver,embedded,mongo,mysql")
		}

		// NOTE : we cannot connect to the docker container instantly after creation. wait for some time before making database connection
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		s.Suffix = fmt.Sprintf("    Waiting for container (%s) to start, it might take about %d minute", dbtype, int(math.Ceil(float64(duration)/60.0)))
		_ = s.Color("green")
		s.Start()
		time.Sleep(time.Duration(duration) * time.Second) // Run for some time to simulate work// Start the spinner
		s.Stop()

		account, err := utils.GetSelectedAccount()
		if err != nil {
			return utils.LogError("Unable to fetch account information", err)
		}
		login, err := utils.Login(account)
		if err != nil {
			return utils.LogError("Unable to login", err)
		}

		v := &model.SpecObject{
			API:  "/v1/config/projects/{project}/database/{dbAlias}/config/{id}",
			Type: "db-config",
			Meta: map[string]string{"project": project, "dbAlias": alias, "id": "database-config"},
			Spec: map[string]interface{}{"conn": connDefault, "type": dbtype, "enabled": true},
		}
		if err := operations.ApplySpec(login.Token, account, v); err != nil {
			utils.LogInfo(`Unable to set database config, try configuring database from mission control`)
			return nil
		}
	}
	utils.LogInfo(fmt.Sprintf("Started database (%s) with alias (%s) & hostname (%s)", dbtype, alias, hostName))
	return nil
}

func removeDatabase(alias string) error {
	ctx := context.Background()

	// Create a docker client
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return utils.LogError("Unable to initialize docker client", err)
	}

	// Check if a database container already exist
	filterArgs := filters.Arg("name", fmt.Sprintf("space-cloud--addon--db--%s", alias))
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filterArgs)})
	if err != nil {
		return utils.LogError("Unable to check if database already exists", err)
	}
	if len(containers) == 0 {
		utils.LogInfo(fmt.Sprintf("Database (%s) not found.", alias))
		return nil
	}

	// Get the hosts file
	hosts, err := txeh.NewHosts(&txeh.HostsConfig{ReadFilePath: utils.GetSpaceCloudHostsFilePath(), WriteFilePath: utils.GetSpaceCloudHostsFilePath()})
	if err != nil {
		return utils.LogError("Unable to open hosts file", err)
	}

	for _, c := range containers {
		hostName := utils.GetServiceDomain("db", alias)

		// Remove the domain from the hosts file
		hosts.RemoveHost(hostName)

		// remove the container from host machine
		if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to remove container %s", c.ID), err)
		}
	}

	// Save the hosts file
	if err := hosts.Save(); err != nil {
		return utils.LogError("Could not save hosts file after updating add on containers", err)
	}

	utils.LogInfo(fmt.Sprintf("Removed database (%s)", alias))

	return nil
}
