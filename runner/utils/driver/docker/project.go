package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/spaceuptech/space-cloud/runner/model"

	"github.com/sirupsen/logrus"
)

func (d *docker) DeleteProject(ctx context.Context, projectID string) error {
	if err := d.DeleteService(ctx, projectID, "", ""); err != nil {
		logrus.Errorf("error deleting project in docker unable to delete services with project id (%s) - %s", projectID, err.Error())
		return err
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/%s", d.secretPath, projectID)); err != nil {
		logrus.Errorf("error deleting project in docker unable to delete secrets folder at (%s) - %s", d.secretPath, err.Error())
		return err
	}
	return nil
}

func (d *docker) CreateProject(ctx context.Context, project *model.Project) error {
	logrus.Debug("create project not implemented for docker")
	return nil
}
