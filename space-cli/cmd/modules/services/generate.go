package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/input"
)

// GenerateService creates a service struct
func GenerateService(projectID, dockerImage string) (*model.SpecObject, error) {
	if projectID == "" {
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID"}, &projectID); err != nil {
			return nil, err
		}
	}

	serviceID := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Service ID"}, &serviceID); err != nil {
		return nil, err
	}

	serviceVersion := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Service Version", Default: "v1"}, &serviceVersion); err != nil {
		return nil, err
	}

	var port int32
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Service Port", Default: "8080"}, &port); err != nil {
		return nil, err
	}

	if dockerImage == "auto" {
		p, err := project.GetProjectConfig(projectID, "project", nil)
		if err != nil {
			return nil, err
		}
		if len(p) == 0 {
			return nil, utils.LogError(fmt.Sprintf("No project found with id (%s)", projectID), err)
		}
		registry, present := p[0].Spec.(map[string]interface{})["dockerRegistry"]
		if registry == "" || !present {
			return nil, fmt.Errorf("no docker registry provided for project (%s)", projectID)
		}

		dockerImage = fmt.Sprintf("%s/%s-%s:%s", registry, projectID, serviceID, serviceVersion)
	}

	want := ""
	dockerSecret := ""
	fileEnvSecret := ""
	secrets := []string{}
	if err := input.Survey.AskOne(&survey.Input{Message: "Are you using private docker registry (Y / N) ?", Default: "N"}, &want); err != nil {
		return nil, err
	}
	if want == "Y" {
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Docker Secret"}, &dockerSecret); err != nil {
			return nil, err
		}
	}

	if err := input.Survey.AskOne(&survey.Input{Message: "Do you want to add other secrets (Y / N) ?", Default: "N"}, &want); err != nil {
		return nil, err
	}
	if want == "Y" {
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter File & Environment Secret (CSV)"}, &fileEnvSecret); err != nil {
			return nil, err
		}
		if fileEnvSecret != "" {
			secrets = append(secrets, strings.Split(fileEnvSecret, ",")...)
		}
	}

	replicaRange := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Replica Range", Default: "1-100"}, &replicaRange); err != nil {
		return nil, err
	}
	replicaMin, replicaMax := 1, 100
	arr := strings.Split(replicaRange, "-")
	if len(arr) != 0 {
		min, err := strconv.Atoi(arr[0])
		if err != nil {
			_ = utils.LogError(fmt.Sprintf("error in generate service config unable to convert replica min which is string to integer - %s", err), nil)
			return nil, err
		}
		replicaMin = min

		max, err := strconv.Atoi(arr[1])
		if err != nil {
			_ = utils.LogError(fmt.Sprintf("error in generate service config unable to convert replica max which is string to integer - %s", err), nil)
			return nil, err
		}
		replicaMax = max
	}

	v := &model.SpecObject{
		API:  "/v1/runner/{project}/services/{id}/{version}",
		Type: "service",
		Meta: map[string]string{
			"id":      serviceID,
			"project": projectID,
			"version": serviceVersion,
		},
		Spec: &model.Service{
			Labels: map[string]string{},
			Scale:  model.ScaleConfig{Replicas: int32(replicaMin), MinReplicas: int32(replicaMin), MaxReplicas: int32(replicaMax), Concurrency: 50, Mode: "parallel"},
			Tasks: []model.Task{
				{
					ID:        serviceID,
					Ports:     []model.Port{{Name: "http", Protocol: "http", Port: port}},
					Resources: model.Resources{CPU: 250, Memory: 512},
					Docker:    model.Docker{ImagePullPolicy: model.PullIfNotExists, Image: dockerImage, Secret: dockerSecret, Cmd: []string{}},
					Runtime:   model.Image,
					Secrets:   secrets,
					Env:       map[string]string{},
				},
			},
			Affinity:  []model.Affinity{},
			Whitelist: []model.Whitelist{{ProjectID: projectID, Service: "*"}},
			Upstreams: []model.Upstream{{ProjectID: projectID, Service: "*"}},
		},
	}
	return v, nil
}

// GenerateServiceRoute creates a service route struct
func GenerateServiceRoute(projectID string) (*model.SpecObject, error) {
	if projectID == "" {
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID:"}, &projectID); err != nil {
			return nil, err
		}
	}

	id := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Route ID:"}, &id); err != nil {
		return nil, err
	}

	var port int32
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Port:", Default: "8080"}, &port); err != nil {
		return nil, err
	}

	version := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Version:"}, &version); err != nil {
		return nil, err
	}

	v := &model.SpecObject{
		API:  "/v1/runner/{project}/service-routes/{id}",
		Type: "service-route",
		Meta: map[string]string{
			"id":      id,
			"project": projectID,
		},
		Spec: map[string]interface{}{
			"routes": []interface{}{
				map[string]interface{}{
					"source": map[string]interface{}{
						"port": port,
					},
					"targets": []interface{}{
						map[string]interface{}{
							"type":    "internal",
							"version": version,
							"weight":  100,
						},
					},
				},
			},
		},
	}

	return v, nil
}

// GenerateServiceRole creates a service role struct
func GenerateServiceRole(projectID string) (*model.SpecObject, error) {
	if projectID == "" {
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Project ID:"}, &projectID); err != nil {
			return nil, err
		}
	}

	roleID := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter Role ID:"}, &roleID); err != nil {
		return nil, err
	}

	serviceID := ""
	if err := input.Survey.AskOne(&survey.Input{Message: "Enter serviceID:"}, &serviceID); err != nil {
		return nil, err
	}
	Types := []string{"Project", "Cluster"}
	Type := ""
	if err := input.Survey.AskOne(&survey.Select{Message: "Choose the Type of role: ", Options: Types, Default: Types[0]}, &Type); err != nil {
		return nil, err
	}

	rules := make([]model.Rule, 0)
	repeat := "YES"
	for repeat == "YES" {
		apiGroup := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter APIGroups:"}, &apiGroup); err != nil {
			return nil, err
		}
		apiGroups := strings.Split(apiGroup, ",")
		if apiGroups[len(apiGroups)-1] == "" {
			return nil, utils.LogError(fmt.Sprintf("Last element of APIGroups not Specified"), nil)
		}

		verb := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Verbs:"}, &verb); err != nil {
			return nil, err
		}
		verbs := strings.Split(verb, ",")
		if verbs[len(verbs)-1] == "" {
			return nil, utils.LogError(fmt.Sprintf("Last element of Verbs not Specified"), nil)
		}

		resource := ""
		if err := input.Survey.AskOne(&survey.Input{Message: "Enter Resources:"}, &resource); err != nil {
			return nil, err
		}
		resources := strings.Split(resource, ",")
		if resources[len(resources)-1] == "" {
			return nil, utils.LogError(fmt.Sprintf("Last element of Resources not Specified"), nil)
		}

		rule := model.Rule{
			APIGroups: apiGroups,
			Verbs:     verbs,
			Resources: resources,
		}
		rules = append(rules, rule)

		again := []string{"YES", "NO"}
		if err := input.Survey.AskOne(&survey.Select{Message: "Do you want to add another rule: ", Options: again, Default: again[1]}, &repeat); err != nil {
			return nil, err
		}
	}
	v := &model.SpecObject{
		API:  "/v1/runner/{project}/service-roles/{serviceId}/{roleId}",
		Type: "service-role",
		Meta: map[string]string{
			"roleId":    roleID,
			"project":   projectID,
			"serviceId": serviceID,
		},
		Spec: &model.Role{
			Type:  Type,
			Rules: rules,
		},
	}
	return v, nil
}
