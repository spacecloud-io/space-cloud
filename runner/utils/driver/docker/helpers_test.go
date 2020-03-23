package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/spaceuptech/space-cloud/runner/model"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type A struct {
}

func (a *A) ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) {
	if imageID == "image-does-not-exist-1" || imageID == "image-does-not-exist-2" {
		return types.ImageInspect{}, nil, fmt.Errorf("provided image doesn't exists")
	}
	return types.ImageInspect{}, nil, nil
}
func (a *A) ImagePull(ctx context.Context, imageID string, options types.ImagePullOptions) (io.ReadCloser, error) {
	if imageID == "image-does-not-exist-1" || imageID == "image-exists" {
		return ioutil.NopCloser(strings.NewReader("pull operation taking place")), nil
	}
	return nil, fmt.Errorf("unable to pull specified image (%s)", imageID)
}
func (a *A) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{}, nil
}
func (a *A) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}
func (a *A) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return types.ContainerJSON{}, nil
}
func (a *A) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return nil, nil
}
func (a *A) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return nil
}

func generateSecretsFile(secret string) error {
	switch secret {
	case "perfect-docker-json-secret-file":
		data := []byte(`{"data":{"username":"ash","password":"123"}}`)
		if err := os.MkdirAll("~secretPath/dummyProject", os.ModePerm); err != nil {
			return err
		}
		return ioutil.WriteFile("~secretPath/dummyProject/perfect-docker-json-secret-file.json", data, os.ModePerm)
	case "imperfect-docker-json-secret-file":
		data := []byte(`"data":{"username":"ash","password":123}`)
		return ioutil.WriteFile("~secretPath/dummyProject/imperfect-docker-json-secret-file.json", data, os.ModePerm)
	}
	return nil
}

func TestDocker_pullImage(t *testing.T) {
	defaultField := &Docker{
		client:       &A{},
		auth:         nil,
		artifactAddr: "",
		secretPath:   "",
		hostFilePath: "",
		manager:      nil,
	}
	type args struct {
		ctx        context.Context
		projectID  string
		taskDocker model.Docker
	}
	tests := []struct {
		name    string
		fields  *Docker
		args    args
		wantErr bool
	}{
		{
			name:   "secrets not required -- image exists",
			fields: defaultField,
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-exists",
				},
			},
			wantErr: false,
		},
		{
			name:   "secrets not required -- image doesn't exists",
			fields: defaultField,
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-exists",
				},
			},
			wantErr: false,
		},
		{
			name: "secrets required -- perfect secret -- image exists",
			fields: &Docker{
				client:     &A{},
				secretPath: "~/secretPath",
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-exists",
					Secret:          "perfect-docker-json-secret-file",
				},
			},
			wantErr: false,
		},
		{
			name: "secrets required -- perfect secret file -- file not present at specified secret path",
			fields: &Docker{
				client:     &A{},
				secretPath: "~misleading-path/secretPath",
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-exists",
					Secret:          "perfect-docker-json-secret-file",
				},
			},
			wantErr: true,
		},
		{
			name: "secrets required -- imperfect secret file -- incorrect json",
			fields: &Docker{
				client:     &A{},
				secretPath: "~secretPath",
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-exists",
					Secret:          "imperfect-docker-json-secret-file",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := generateSecretsFile(tt.args.taskDocker.Secret); err != nil {
				t.Errorf("unable to create docker secret file which is a dependency for running test succesfully (%s)", err.Error())
			}

			if err := tt.fields.pullImage(tt.args.ctx, tt.args.projectID, tt.args.taskDocker); (err != nil) != tt.wantErr {
				t.Errorf("pullImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	os.RemoveAll("~secretPath")
}

func TestDocker_pullImageByPolicy(t *testing.T) {
	type args struct {
		ctx        context.Context
		projectID  string
		taskDocker model.Docker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "image policy pull if not exits -- image exists",
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-exists",
				},
			},
			wantErr: false,
		},
		{
			name: "image policy pull if not exists -- image doesn't exists -- successful image pull",
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-does-not-exist-1",
					Secret:          "",
				},
			},
			wantErr: false,
		},
		{
			name: "image policy pull if not exists -- image doesn't exists -- unsuccessful image pull",
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "image-does-not-exist-2",
					Secret:          "",
				},
			},
			wantErr: true,
		},
		{
			name: "image policy always pull -- image doesn't exists -- successful image pull",
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullAlways,
					Image:           "image-does-not-exist-1",
					Secret:          "",
				},
			},
			wantErr: false,
		},
		{
			name: "image policy always pull -- image doesn't exists -- unsuccessful image pull",
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullAlways,
					Image:           "image-does-not-exist-2",
					Secret:          "",
				},
			},
			wantErr: true,
		},
	}

	d := &Docker{
		client:       &A{},
		auth:         nil,
		artifactAddr: "",
		secretPath:   "",
		hostFilePath: "",
		manager:      nil,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.pullImageByPolicy(tt.args.ctx, tt.args.projectID, tt.args.taskDocker); (err != nil) != tt.wantErr {
				t.Errorf("pullImageByPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
