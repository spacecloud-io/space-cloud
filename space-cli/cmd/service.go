package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/model"
)

// GenerateService creates a service struct
func GenerateService() (*model.Service, error) {
	account, err := getSelectedAccount()
	if err != nil {
		logrus.Errorf("error in generate service unable to get selected account - %s", err.Error())
		return nil, err
	}

	loginResult, err := login(account)
	if err != nil {
		logrus.Errorf("error in generate service unable to login - %s", err.Error())
		return nil, err
	}

	// create new services.yaml file
	return generateServiceConfig(loginResult.Projects)
}

func generateServiceConfig(projects []*model.Projects) (*model.Service, error) {
	serviceID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Service ID"}, &serviceID); err != nil {
		return nil, err
	}

	serviceVersion := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Service Version"}, &serviceVersion); err != nil {
		return nil, err
	}
	var port int32
	if err := survey.AskOne(&survey.Input{Message: "Enter Service Port"}, &port); err != nil {
		return nil, err
	}
	projectID := ""
	projectNames, err := getProjects(projects)
	if err != nil {
		return nil, err
	}
	if err := survey.AskOne(&survey.Select{Message: "Select Project", Options: projectNames}, &projectID); err != nil {
		return nil, err
	}
	dockerImage := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Docker Image Name"}, &dockerImage); err != nil {
		return nil, err
	}
	if dockerImage == "" {
		dockerImage = fmt.Sprintf("%s/%s", projectID, serviceID)
	}

	want := ""
	dockerSecret := ""
	fileEnvSecret := ""
	secrets := []string{}
	if err := survey.AskOne(&survey.Input{Message: "Is this a private repository (Y / N) ?", Default: "N"}, &want); err != nil {
		return nil, err
	}
	if want == "Y" {
		if err := survey.AskOne(&survey.Input{Message: "Enter Docker Secret"}, &dockerSecret); err != nil {
			return nil, err
		}
	}

	if err := survey.AskOne(&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want); err != nil {
		return nil, err
	}
	if want == "Y" {
		if err := survey.AskOne(&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &fileEnvSecret); err != nil {
			return nil, err
		}
		if fileEnvSecret != "" {
			secrets = append(secrets, strings.Split(fileEnvSecret, ",")...)
		}
	}

	replicaRange := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Replica Range", Default: "0-100"}, &replicaRange); err != nil {
		return nil, err
	}
	replicaMin, replicaMax := 0, 100
	arr := strings.Split(replicaRange, "-")
	if len(arr) != 0 {
		min, err := strconv.Atoi(arr[0])
		if err != nil {
			logrus.Errorf("error in generate service config unable to convert replica min which is string to integer - %s", err)
			return nil, err
		}
		replicaMin = min

		max, err := strconv.Atoi(arr[1])
		if err != nil {
			logrus.Errorf("error in generate service config unable to convert replica max which is string to integer - %s", err)
			return nil, err
		}
		replicaMax = max
	}

	c := &model.Service{
		ID:        serviceID,
		ProjectID: projectID,
		Version:   serviceVersion,
		Labels:    map[string]string{},
		Scale:     model.ScaleConfig{Replicas: 0, MinReplicas: int32(replicaMin), MaxReplicas: int32(replicaMax), Concurrency: 50},
		Tasks: []model.Task{
			{
				ID:        serviceID,
				Ports:     []model.Port{{Name: "http", Protocol: "http", Port: port}},
				Resources: model.Resources{CPU: 250, Memory: 512},
				Docker:    model.Docker{Image: dockerImage, Secret: dockerSecret, Cmd: []string{}},
				Runtime:   model.Image,
				Secrets:   secrets,
				Env:       map[string]string{},
			},
		},
		Affinity:  []model.Affinity{},
		Whitelist: []model.Whitelist{{ProjectID: projectID, Service: "*"}},
		Upstreams: []model.Upstream{{ProjectID: projectID, Service: "*"}},
	}
	return c, nil
}
