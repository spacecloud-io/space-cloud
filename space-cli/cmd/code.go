package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/spaceuptech/space-cli/model"
)

// CodeStart starts the code commands
func CodeStart() (*model.Service, *model.LoginResponse, error) {
	credential, err := getCreds()
	if err != nil {
		return nil, nil, err
	}

	selectedAccount := getSelectedAccount(credential)

	loginRes, err := login(selectedAccount)
	if err != nil {
		return nil, nil, err
	}

	c, err := getServiceConfig("services.yaml")
	if err != nil {
		c, err = GenerateServiceConfig(loginRes.Projects, selectedAccount)
		if err != nil {
			return nil, nil, err
		}
	}
	return c, loginRes, nil
}

func GenerateServiceConfig(projects []*model.Projects, selectedaccount *model.Account) (*model.Service, error) {
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
	// todo ask secrets and set in tasks
	c := &model.Service{
		ID:        serviceID,
		Name:      serviceName,
		ProjectID: projectNameID,
		Version:   "v1",
		Scale:     model.ScaleConfig{Replicas: 0, MinReplicas: 0, MaxReplicas: 100, Concurrency: 50},
		Tasks: []model.Task{
			{
				ID:        serviceID,
				Name:      serviceName,
				Ports:     []model.Port{{Name: "http", Protocol: "http", Port: port}},
				Resources: model.Resources{CPU: 250, Memory: 512},
				Docker:    model.Docker{Image: img},
				Env:       map[string]string{"URL": selectedaccount.ServerUrl, "CMD": progCmd}, // todo check why do we need this envs
				Runtime:   model.Image,
			},
		},
		Whitelist: []model.Whitelist{{ProjectID: projectNameID, Service: "*"}},
		Upstreams: []model.Upstream{{ProjectID: projectNameID, Service: "*"}},
	}
	return c, nil
}

// RunDockerFile starts a container using go docker client
func RunDockerFile(s *model.ActionCode, loginResp *model.LoginResponse) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	// sa, err := json.Marshal(s)
	// if err != nil {
	// 	return err
	// }
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/Dockerfile", dir))
	if err != nil {
		return err
	}
	_, err = cli.ImageBuild(ctx, strings.NewReader(string(data)), types.ImageBuildOptions{})
	if err != nil {
		return err
	}
	// todo store some name or for creating container
	return nil
}

func GenerateServiceConfigWithoutLogin() (*model.Service, error) {
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
	if err := survey.AskOne(&survey.Input{Message: "Enter Project Name"}, &projectNameID); err != nil {
		return nil, err
	}

	dockerImage := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Docker Image Name"}, &dockerImage); err != nil {
		return nil, err
	}
	if dockerImage == "" {
		dockerImage = fmt.Sprintf("%s/%s", projectNameID, serviceID)
	}

	dockerSecret := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Docker Secret"}, &dockerSecret); err != nil {
		return nil, err
	}

	secret := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &secret); err != nil {
		return nil, err
	}
	secrets := []string{}
	for _, value := range strings.Split(secret, ",") {
		secrets = append(secrets, value)
	}

	var progCmd string
	if err := survey.AskOne(&survey.Input{Message: "Enter Run Cmd", Default: strings.Join(getCmd(progLang), " ")}, &progCmd); err != nil {
		return nil, err
	}

	c := &model.Service{
		ID:        serviceID,
		Name:      serviceName,
		ProjectID: projectNameID,
		Version:   "v1",
		Scale:     model.ScaleConfig{Replicas: 0, MinReplicas: 0, MaxReplicas: 100, Concurrency: 50},
		Tasks: []model.Task{
			{
				ID:        serviceID,
				Name:      serviceName,
				Ports:     []model.Port{{Name: "http", Protocol: "http", Port: port}},
				Resources: model.Resources{CPU: 250, Memory: 512},
				Docker:    model.Docker{Image: dockerImage, Secret: dockerSecret},
				Env:       map[string]string{}, // todo check why do we need this envs
				Runtime:   model.Image,
				Secrets:   secrets,
			},
		},
		Whitelist: []model.Whitelist{{ProjectID: projectNameID, Service: "*"}},
		Upstreams: []model.Upstream{{ProjectID: projectNameID, Service: "*"}},
	}
	return c, nil
}
