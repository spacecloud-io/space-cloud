package utils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/filters"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
)

// GetContainers gets container running in a specific cluster
func GetContainers(ctx context.Context, dockerClient *client.Client, clusterName, containerType string) ([]types.Container, error) {
	argsName := filters.Arg("name", "space-cloud")
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(argsName), All: true})
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to list containers - %s", err.Error()), nil)
		return nil, err
	}

	clusterContainers := make([]types.Container, 0)
	for _, container := range containers {
		bigArr := strings.Split(container.Names[0], "--")
		smallArr := strings.Split(bigArr[0], "-")
		if clusterName == "default" {
			if len(smallArr) == 3 || len(smallArr) == 2 {
				if containerType == model.DbContainers && len(smallArr) != 2 {
					continue
				}
				// Containers running in a default cluster
				if isTypeSpecificContainer(containerType, container.Labels) {
					clusterContainers = append(clusterContainers, container)
				}
			}
			continue
		}
		if len(smallArr) >= 3 && smallArr[2] == clusterName {
			if isTypeSpecificContainer(containerType, container.Labels) {
				// Containers running in a specific cluster
				clusterContainers = append(clusterContainers, container)
			}
		}
	}
	return clusterContainers, nil
}

func isTypeSpecificContainer(cType string, labels map[string]string) bool {
	switch cType {
	case model.AllContainers:
		return true
	case model.DbContainers:
		return labels["service"] == "db"
	case model.ServiceContainers:
		_, ok1 := labels["service"]
		_, ok2 := labels["task"]
		_, ok3 := labels["version"]
		return ok1 && ok2 && ok3
	case model.ScContainers:
		return labels["service"] == "gateway" || labels["service"] == "runner"
	case model.RegistryContainers:
		return labels["service"] == "registry"
	default:
		return false
	}
}

// PullImageIfNotExist pulls the docker image if it does not exist
func PullImageIfNotExist(ctx context.Context, dockerClient *client.Client, image string) error {
	_, _, err := dockerClient.ImageInspectWithRaw(ctx, image)
	if err != nil {
		// pull image from public repository
		logrus.Infof("Image %s does not exist. Need to pull it from Docker Hub. This may take some time.", image)
		out, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			logrus.Errorf("Unable to pull public image with id (%s) - %s", image, err.Error())
			return err
		}
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		s.Suffix = "    Downloading image..."
		_ = s.Color("green")
		s.Start()
		time.Sleep(4 * time.Second) // Run for some time to simulate work// Start the spinner
		_, _ = io.Copy(ioutil.Discard, out)
		s.Stop()
	}
	logrus.Infof("Image %s already exists. No need to pull it again", image)
	return nil
}

// DockerfileGolang is the docker file to use for golang projects
const DockerfileGolang string = `
FROM golang:1.13.5-alpine3.10
WORKDIR /build
COPY . .
#RUN apk --no-cache add build-base
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

FROM alpine:3.10
WORKDIR /app
COPY --from=0 /build/app .
CMD ["./app"]
`

// DockerfileNodejs is the docker file to use for node js projects
const DockerfileNodejs string = `
FROM node:10-alpine

# Create app directory
WORKDIR /app

COPY package*.json ./

RUN npm install --only=prod

# Bundle app source
COPY . .

CMD [ "node", "index.js" ]`

// DockerfilePython is the docker file to use for python projects
const DockerfilePython string = `
FROM python:3.6

# Create app directory
WORKDIR /app

# Install app dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Bundle app source
COPY . .

CMD [ "python", "app.py" ]`
