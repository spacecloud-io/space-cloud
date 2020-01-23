package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"github.com/spaceuptech/space-cli/model"
)

// CodeStart starts the code commands
func CodeStart(envID string) (*model.Service, *model.LoginResponse, error) {
	credential, err := getCreds()
	if err != nil {
		return nil, nil, err
	}

	selectedAccount := getSelectedAccount(credential)

	loginRes, err := login(selectedAccount)
	if err != nil {
		return nil, nil, err
	}

	c, err := getServiceConfig(getAccountConfigPath())
	if err != nil {
		c, err = generateServiceConfig(loginRes.Projects, selectedAccount, envID)
		if err != nil {
			return nil, nil, err
		}
	}
	return c, loginRes, nil
}

func generateServiceConfig(projects []*model.Projects, selectedaccount *model.Account, envID string) (*model.Service, error) {
	progLang, err := getProgLang()
	if err != nil {
		return nil, err
	}
	serviceName := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Service Name"}, &serviceName); err != nil {
		return nil, err
	}
	defaultServiceID := strings.ReplaceAll(serviceName, " ", "-")
	serviceID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Service ID", Default: strings.ToLower(defaultServiceID)}, &serviceID); err != nil {
		return nil, err
	}
	var port int32
	if err := survey.AskOne(&survey.Input{Message: "Enter Service Port"}, &port); err != nil {
		return nil, err
	}
	projectNameID := ""
	projs, err := getProjects(projects)
	if err != nil {
		return nil, err
	}
	if err := survey.AskOne(&survey.Select{Message: "Select Project Name", Options: projs}, &projectNameID); err != nil {
		return nil, err
	}

	var progCmd string
	if err := survey.AskOne(&survey.Input{Message: "Enter Run Cmd", Default: strings.Join(getCmd(progLang), " ")}, &progCmd); err != nil {
		return nil, err
	}
	img, err := getImage(progLang)
	if err != nil {
		return nil, err
	}

	c := &model.Service{
		ID:          serviceID,
		Name:        serviceName,
		ProjectID:   projectNameID,
		Environment: envID,
		Version:     "v1",
		Scale:       model.ScaleConfig{Replicas: 0, MinReplicas: 0, MaxReplicas: 100, Concurrency: 50},
		Tasks: []model.Task{
			{
				ID:        serviceID,
				Name:      serviceName,
				Ports:     []model.Port{model.Port{Protocol: "http", Port: port}},
				Resources: model.Resources{CPU: 250, Memory: 512},
				Docker:    model.Docker{Image: img},
				Env:       map[string]string{"URL": selectedaccount.ServerUrl, "CMD": progCmd},
			},
		},
		Whitelist: []string{"project:*"},
		Upstreams: []model.Upstream{model.Upstream{ProjectID: projectNameID, Service: "*"}},
		Runtime:   "code",
	}
	return c, nil
}

// RunDockerFile starts a container using go docker client
func RunDockerFile(s *model.ActionCode, loginResp *model.LoginResponse) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	sa, err := json.Marshal(s)
	if err != nil {
		return err
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: s.Service.Tasks[0].Docker.Image,
			Env: []string{
				"FILE_PATH=/",
				fmt.Sprintf("URL=%s", s.Service.Tasks[0].Env["URL"]),
				fmt.Sprintf("TOKEN=%s", loginResp.Token),
				fmt.Sprintf("meta=%s", string(sa))},
		},
		&container.HostConfig{Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: dir,
				Target: "/build",
			},
		}}, nil, "")
	if err != nil {
		return err
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}
