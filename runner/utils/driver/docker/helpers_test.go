package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDocker_pullImage(t *testing.T) {

	type mockArg struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	type args struct {
		ctx        context.Context
		projectID  string
		taskDocker model.Docker
	}
	tests := []struct {
		name                 string
		dockerClientMockArgs []mockArg
		fileSystemMockArgs   []mockArg
		secretPath           string
		args                 args
		wantErr              bool
	}{
		{
			name: "secrets not required -- image exists",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image", types.ImagePullOptions{}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("pull operation taking place")), nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					Image: "unknown-image",
				},
			},
			wantErr: false,
		},
		{
			name: "secrets not required -- image doesn't exists",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image-1", types.ImagePullOptions{}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("error downloading image")), fmt.Errorf("unable to download specified image")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					Image: "unknown-image-1",
				},
			},
			wantErr: true,
		},
		{
			name: "secrets required -- perfect secret -- image exists",
			fileSystemMockArgs: []mockArg{
				{
					method:        "ReadSecretsFiles",
					args:          []interface{}{context.Background(), "dummyProject", "dummySecretFileName"},
					paramReturned: []interface{}{[]byte(`{"data":{"username":"admin","password":"123"}}`), nil},
				},
			},
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image-2", types.ImagePullOptions{RegistryAuth: "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIxMjMifQ=="}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("pull operation taking place")), nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					Image:  "unknown-image-2",
					Secret: "dummySecretFileName",
				},
			},
			wantErr: false,
		},
		{
			name: "secrets required -- perfect secret file -- file not present at specified secret path",
			fileSystemMockArgs: []mockArg{
				{
					method:        "ReadSecretsFiles",
					args:          []interface{}{context.Background(), "dummyProject", "dummySecretFileName-1"},
					paramReturned: []interface{}{[]byte(""), fmt.Errorf("unable to open file at specified path")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					Image:  "unknown-image-2",
					Secret: "dummySecretFileName-1",
				},
			},
			wantErr: true,
		},
		{
			name: "secrets required -- imperfect secret file -- incorrect json",
			fileSystemMockArgs: []mockArg{
				{
					method:        "ReadSecretsFiles",
					args:          []interface{}{context.Background(), "dummyProject", "dummySecretFileName-2"},
					paramReturned: []interface{}{[]byte(`{data":{"username":"admin","password":"123"}}`), nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					Image:  "unknown-image-2",
					Secret: "dummySecretFileName-2",
				},
			},
			wantErr: true,
		},
	}

	t.Parallel()
	assertObj := assert.New(t)
	a := &dockerClientMock{}
	f := &fileSystemMock{}
	d := &Docker{
		client:       a,
		auth:         nil,
		artifactAddr: "",
		secretPath:   "",
		hostFilePath: "",
		manager:      nil,
		fileSystem:   f,
	}

	for _, tt := range tests {
		d.secretPath = tt.secretPath
		for _, value := range tt.dockerClientMockArgs {
			a.On(value.method, value.args...).Return(value.paramReturned...)
		}
		for _, value := range tt.fileSystemMockArgs {
			f.On(value.method, value.args...).Return(value.paramReturned...)
		}
		err := d.pullImage(tt.args.ctx, tt.args.projectID, tt.args.taskDocker)
		assertObj.Equal(tt.wantErr, err != nil, fmt.Sprintf("Test name (%s) Got error (%v)", tt.name, err))
	}
}

func TestDocker_pullImageByPolicy(t *testing.T) {
	type mockArg struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	type args struct {
		ctx        context.Context
		projectID  string
		taskDocker model.Docker
	}
	tests := []struct {
		name                 string
		dockerClientMockArgs []mockArg
		args                 args
		wantErr              bool
	}{
		{
			name: "image policy pull if not exits -- image exists",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImageInspectWithRaw",
					args:          []interface{}{context.Background(), "dummyImage"},
					paramReturned: []interface{}{types.ImageInspect{}, []byte("success"), nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "dummyImage",
				},
			},
			wantErr: false,
		},
		{
			name: "image policy pull if not exists -- image doesn't exists -- successful image pull",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImageInspectWithRaw",
					args:          []interface{}{context.Background(), "unknown-image"},
					paramReturned: []interface{}{types.ImageInspect{}, []byte("unsuccessful"), fmt.Errorf("unable to find specifed image")},
				},
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image", types.ImagePullOptions{}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("pull operation taking place")), nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "unknown-image",
					Secret:          "",
				},
			},
			wantErr: false,
		},
		{
			name: "image policy pull if not exists -- image doesn't exists -- unsuccessful image pull",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImageInspectWithRaw",
					args:          []interface{}{context.Background(), "unknown-image-1"},
					paramReturned: []interface{}{types.ImageInspect{}, []byte("unsuccessful"), fmt.Errorf("unable to find specifed image")},
				},
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image-1", types.ImagePullOptions{}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("error downloading image")), fmt.Errorf("unable to pull image from docker repository")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullIfNotExists,
					Image:           "unknown-image-1",
					Secret:          "",
				},
			},
			wantErr: true,
		},
		{
			name: "image policy always pull -- image doesn't exists -- successful image pull",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImageInspectWithRaw",
					args:          []interface{}{context.Background(), "unknown-image"},
					paramReturned: []interface{}{types.ImageInspect{}, []byte("unsuccessful"), fmt.Errorf("unable to find specifed image")},
				},
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image", types.ImagePullOptions{}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("pull operation taking place")), nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullAlways,
					Image:           "unknown-image",
					Secret:          "",
				},
			},
			wantErr: false,
		},
		{
			name: "image policy always pull -- image doesn't exists -- unsuccessful image pull",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ImageInspectWithRaw",
					args:          []interface{}{context.Background(), "unknown-image-1"},
					paramReturned: []interface{}{types.ImageInspect{}, []byte("unsuccessful"), fmt.Errorf("unable to find specifed image")},
				},
				{
					method:        "ImagePull",
					args:          []interface{}{context.Background(), "unknown-image-1", types.ImagePullOptions{}},
					paramReturned: []interface{}{ioutil.NopCloser(strings.NewReader("error downloading image")), fmt.Errorf("unable to pull image from docker repository")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject",
				taskDocker: model.Docker{
					ImagePullPolicy: model.PullAlways,
					Image:           "unknown-image-1",
					Secret:          "",
				},
			},
			wantErr: true,
		},
	}

	t.Parallel()
	assertObj := assert.New(t)
	a := &dockerClientMock{}
	d := &Docker{
		client:       a,
		auth:         nil,
		artifactAddr: "",
		secretPath:   "",
		hostFilePath: "",
		manager:      nil,
	}

	for _, tt := range tests {
		for _, value := range tt.dockerClientMockArgs {
			a.On(value.method, value.args...).Return(value.paramReturned...)
		}
		assertObj.Equal(tt.wantErr, d.pullImageByPolicy(tt.args.ctx, tt.args.projectID, tt.args.taskDocker) != nil, tt.name)
		// a.AssertExpectations(t)
	}
}
