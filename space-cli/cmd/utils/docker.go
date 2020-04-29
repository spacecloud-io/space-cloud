package utils

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/briandowns/spinner"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

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
