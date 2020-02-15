package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func getServiceDomain(projectID, serviceID string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", serviceID, projectID)
}

func (d *docker) pullImageIfDoesntExists(ctx context.Context, projectId string, taskDocker model.Docker) error {
	_, _, err := d.client.ImageInspectWithRaw(ctx, taskDocker.Image)
	if err != nil {
		// image doesn't exist locally
		if taskDocker.Secret != "" {
			data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s.json", d.secretPath, projectId, taskDocker.Secret))
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

			// pull image from private repository
			out, err := d.client.ImagePull(ctx, taskDocker.Image, types.ImagePullOptions{RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON)})
			if err != nil {
				logrus.Errorf("error in docker unable to pull private image with id (%s) - %s", taskDocker.Image, err.Error())
				return err
			}
			io.Copy(os.Stdout, out)
		} else {
			// pull image from public repository
			out, err := d.client.ImagePull(ctx, taskDocker.Image, types.ImagePullOptions{})
			if err != nil {
				logrus.Errorf("error in docker unable to pull public image with id (%s) - %s", taskDocker.Image, err.Error())
				return err
			}
			io.Copy(os.Stdout, out)
		}
	}
	return nil
}

func getSpaceCloudHostsFilePath() string {
	return fmt.Sprintf("%s/hosts", getSpaceCloudDirectory())
}

func getSpaceCloudDirectory() string {
	return fmt.Sprintf("%s/.space-cloud", getHomeDirectory())
}

func getHomeDirectory() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
