package addons

import (
	"context"
	"fmt"
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
			return utils.LogError(`Cannot fetch quotas. Is Space Cloud running?`, err)
		}

		// fetch current db config
		dbConfig, err := database.GetDbConfig(project, "db-config", map[string]string{})
		if err != nil {
			return utils.LogError(`Cannot fetch database config from gateway`, err)
		}

		// check if database can be added
		if len(dbConfig) >= resp.Quotas.MaxDatabases {
			return utils.LogError(fmt.Sprintf("Cannot add database in project (%s), max database limit reached. upgrade you plan", project), err)
		}
	}

	// Prepare the docker image name name
	dockerImage := fmt.Sprintf("%s:%s", dbtype, version)
	if dbtype == "sqlserver" {
		dockerImage = fmt.Sprintf("%s:%s", "mcr.microsoft.com/mssql/server", version)
	}
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
	case "sqlserver":
		env = []string{"ACCEPT_EULA=Y", fmt.Sprintf("SA_PASSWORD=%s", password)}
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

	utils.LogInfo(fmt.Sprintf("Started database (%s) with alias (%s) & hostname (%s)", dbtype, alias, hostName))

	if autoApply {
		connDefault := ""
		switch dbtype {
		case "postgres":
			connDefault = fmt.Sprintf("postgres://%s:%s@%s:5432/postgres?sslmode=disable", username, password, hostName)
		case "mongo":
			connDefault = fmt.Sprintf("mongodb://%s:27017", hostName)
		case "mysql":
			connDefault = fmt.Sprintf("root:%s@tcp(%s:3306)/", password, hostName)
		case "sqlserver":
			connDefault = fmt.Sprintf("Data Source=%s,1433;Initial Catalog=master;User ID=%s;Password=%s;", hostName, username, password)
		default:
			return fmt.Errorf("invalid database provided, supported databases postgres,sqlserver,embedded,mongo,mysql")
		}

		account, token, err := utils.LoginWithSelectedAccount()
		if err != nil {
			return utils.LogError("Couldn't get account details or login token", err)
		}

		v := &model.SpecObject{
			API:  "/v1/config/projects/{project}/database/{dbAlias}/config/{id}",
			Type: "db-config",
			Meta: map[string]string{"project": project, "dbAlias": alias, "id": "database-config"},
			Spec: map[string]interface{}{"conn": connDefault, "type": dbtype, "enabled": true},
		}
		keepSettingConfig(token, dbtype, account, v)

	}
	return nil
}

func keepSettingConfig(token, dbType string, account *model.Account, v *model.SpecObject) {
	timeout := time.After(5 * time.Minute) // 5 is the maximum time required as mysql may take upto 5 minutes
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	// NOTE : we cannot connect to the docker container instantly after creation. wait for some time before making database connection
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	s.Suffix = fmt.Sprintf("    Waiting for container (%s) to start", dbType)
	_ = s.Color("green")
	s.Start()
	defer s.Stop()
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			logrus.Warningln(`Unable to set database config, try configuring database from mission control`)
			return
			// Got a tick, we should check on checkSomething()
		case <-ticker.C:
			if err := operations.ApplySpec(token, account, v); err != nil {
				logrus.Warningln("Unable to add database to Space Cloud config", nil)
				continue
			}

			v = &model.SpecObject{
				API:  "/v1/config/projects/{project}/database/{dbAlias}/collections/{col}/rules",
				Type: "db-rules",
				Meta: map[string]string{
					"dbAlias": v.Meta["dbAlias"],
					"col":     "default",
					"project": v.Meta["project"],
				},
				Spec: map[string]interface{}{
					"isRealtimeEnabled": false,
					"rules": map[string]interface{}{
						"create": map[string]interface{}{
							"rule": "allow",
						},
						"delete": map[string]interface{}{
							"rule": "allow",
						},
						"read": map[string]interface{}{
							"rule": "allow",
						},
						"update": map[string]interface{}{
							"rule": "allow",
						},
					},
				},
			}
			if err := operations.ApplySpec(token, account, v); err != nil {
				logrus.Warningln("Couldn't add default collection rules", nil)
			}

			v = &model.SpecObject{
				API:  "/v1/config/projects/{project}/database/{dbAlias}/prepared-queries/{id}",
				Type: "eventing-rule",
				Meta: map[string]string{
					"project": v.Meta["project"],
					"dbAlias": v.Meta["dbAlias"],
					"id":      "default",
				},
				Spec: map[string]interface{}{
					"id":  "default",
					"sql": "",
					"rule": map[string]interface{}{
						"rule": "allow",
					},
					"args": []string{},
				},
			}
			if err := operations.ApplySpec(token, account, v); err != nil {
				logrus.Warningln("Couldn't add default prepared query rules", nil)
			}

			utils.LogInfo("Successfully added database to Space Cloud config.")
			return
		}
	}
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
