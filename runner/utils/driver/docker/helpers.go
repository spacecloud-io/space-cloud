package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func (d *Docker) pullImageByPolicy(ctx context.Context, projectID string, taskDocker model.Docker) error {
	if taskDocker.ImagePullPolicy == model.PullIfNotExists {
		if _, _, err := d.client.ImageInspectWithRaw(ctx, taskDocker.Image); err == nil {
			return nil
		}
	}
	err := d.pullImage(ctx, projectID, taskDocker)
	if err != nil {
		return err
	}
	return nil
}

func (d *Docker) pullImage(ctx context.Context, projectID string, taskDocker model.Docker) error {
	// image doesn't exist locally
	options := types.ImagePullOptions{}
	if taskDocker.Secret != "" {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectID, taskDocker.Secret))
		if err != nil {
			logrus.Errorf("error in docker unable to read file name (%s) required for pulling image from private repository - %v", taskDocker.Secret, err.Error())
			return err
		}
		secret := new(model.Secret)
		if err := json.Unmarshal(data, secret); err != nil {
			return err
		}

		authConfig := types.AuthConfig{
			Username: secret.Data["username"],
			Password: secret.Data["password"],
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}

		options.RegistryAuth = base64.URLEncoding.EncodeToString(encodedJSON)
	}

	// pull image from repository
	out, err := d.client.ImagePull(ctx, taskDocker.Image, options)
	if err != nil {
		logrus.Errorf("error in docker unable to pull private image with id (%s) - %s", taskDocker.Image, err.Error())
		return err
	}
	_, _ = io.Copy(ioutil.Discard, out)
	return nil
}
