package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/stretchr/testify/assert"
	"github.com/txn2/txeh"
	"reflect"
	"testing"
)

func TestDocker_ApplyService(t *testing.T) {
	// type mockArg struct {
	// 	method        string
	// 	args          []interface{}
	// 	paramReturned []interface{}
	// }
	// type args struct {
	// 	ctx     context.Context
	// 	service *model.Service
	// }
	// tests := []struct {
	// 	name                 string
	// 	dockerClientMockArgs []mockArg
	// 	fileSystemMockArgs   []mockArg
	// 	proxyManagerMockArgs []mockArg
	// 	args                 args
	// 	wantErr              bool
	// }{
	// 	{
	// dockerClientMockArgs: []mockArg{
	// 	{
	// 		method:        "ContainerList",
	// 		args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject1--dummyService1--dummyVersion1")), All: true}},
	// 		paramReturned: []interface{}{[]types.Container{{ID: "dummyService1-container"}}, nil},
	// 	},
	// 	{
	// 		method:        "ContainerRemove",
	// 		args:          []interface{}{context.Background(), "dummyService1-container", types.ContainerRemoveOptions{Force: true}},
	// 		paramReturned: []interface{}{nil},
	// 	},
	// 	{
	// 		method:        "ContainerList",
	// 		args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject1--dummyService1")), All: true}},
	// 		paramReturned: []interface{}{[]types.Container{{ID: "some-container-exists"}}, nil},
	// 	},
	// },
	// fileSystemMockArgs: []mockArg{
	// 	{
	// 		method:        "NewHostFile",
	// 		args:          []interface{}{},
	// 		paramReturned: []interface{}{&txeh.Hosts{}, nil},
	// 	},
	//
	// 	{
	// 		method:        "NewHostFile",
	// 		args:          []interface{}{},
	// 		paramReturned: []interface{}{&txeh.Hosts{}, nil},
	// 	},
	// 	{
	// 		method:        "RemoveTempSecretsFolder",
	// 		args:          []interface{}{"dummyProject1", "dummyService1", "dummyVersion1"},
	// 		paramReturned: []interface{}{nil},
	// 	},
	// 	{
	// 		method:        "RemoveHostFromHostFile",
	// 		args:          []interface{}{&txeh.Hosts{}, utils.GetInternalServiceDomain("dummyProject1", "dummyService1", "dummyVersion1")},
	// 		paramReturned: []interface{}{},
	// 	},
	// 	{
	// 		method:        "SaveHostFile",
	// 		args:          []interface{}{&txeh.Hosts{}},
	// 		paramReturned: []interface{}{nil},
	// 	},
	// },
	// args: args{
	// 	ctx: context.Background(),
	// 	service: &model.Service{
	// 		ID:        "dummyService1",
	// 		Name:      "dummyService1Name1",
	// 		ProjectID: "dummyProject1",
	// 		Version:   "dummyVersion1",
	// 		Scale:     model.ScaleConfig{},
	// 		Labels:    nil,
	// 		Tasks: []model.Task{
	// 			{
	// 				ID:   "dummyTask1",
	// 				Name: "dummyTask1Name1",
	// 				Ports: []model.Port{
	// 					{
	// 						Name:     "dummyProtocol1",
	// 						Protocol: "http",
	// 						Port:     8080,
	// 					},
	// 				},
	// 				Resources: model.Resources{},
	// 				Docker:    model.Docker{},
	// 				Env:       nil,
	// 				Secrets:   nil,
	// 				Runtime:   "",
	// 			},
	// 			{
	// 				ID:   "dummyTask2",
	// 				Name: "dummyTask2Name2",
	// 				Ports: []model.Port{
	// 					{
	// 						Name:     "dummyProtocol2",
	// 						Protocol: "http",
	// 						Port:     8081,
	// 					},
	// 				},
	// 				Resources: model.Resources{},
	// 				Docker:    model.Docker{},
	// 				Env:       nil,
	// 				Secrets:   nil,
	// 				Runtime:   "",
	// 			},
	// 		},
	// 		Affinity:  nil,
	// 		Whitelist: nil,
	// 		Upstreams: nil,
	// 	},
	// },
	// 	},
	// }
	// t.Parallel()
	// assertObj := assert.New(t)
	// a := &dockerClientMock{}
	// f := &fileSystemMock{}
	// p := &proxyMangerMock{}
	// d := &Docker{
	// 	client:       a,
	// 	auth:         nil,
	// 	artifactAddr: "",
	// 	secretPath:   "",
	// 	hostFilePath: "",
	// 	manager:      p,
	// 	fileSystem:   f,
	// }
	//
	// for _, tt := range tests {
	// 	for _, value := range tt.dockerClientMockArgs {
	// 		a.On(value.method, value.args...).Return(value.paramReturned...)
	// 	}
	// 	for _, value := range tt.fileSystemMockArgs {
	// 		f.On(value.method, value.args...).Return(value.paramReturned...)
	// 	}
	// 	for _, value := range tt.proxyManagerMockArgs {
	// 		p.On(value.method, value.args...).Return(value.paramReturned...)
	// 	}
	// 	err := d.ApplyService(tt.args.ctx, tt.args.service)
	// 	assertObj.Equal(tt.wantErr, err != nil, fmt.Sprintf("Test name (%s) Got error (%v)", tt.name, err))
	// }
	// a.AssertExpectations(t)
	// f.AssertExpectations(t)
	// p.AssertExpectations(t)
}

func TestDocker_ApplyServiceRoutes(t *testing.T) {
	// type mockArg struct {
	// 	method        string
	// 	args          []interface{}
	// 	paramReturned []interface{}
	// }
	// type args struct {
	// 	ctx       context.Context
	// 	projectID string
	// 	serviceID string
	// 	routes    model.Routes
	// }
	// tests := []struct {
	// 	name                 string
	// 	dockerClientMockArgs []mockArg
	// 	fileSystemMockArgs   []mockArg
	// 	proxyManagerMockArgs []mockArg
	// 	args                 args
	// 	wantErr              bool
	// }{
	// 	{
	//
	// 	},
	// }
	//
	// t.Parallel()
	// assertObj := assert.New(t)
	// a := &dockerClientMock{}
	// f := &fileSystemMock{}
	// p := &proxyMangerMock{}
	// d := &Docker{
	// 	client:       a,
	// 	auth:         nil,
	// 	artifactAddr: "",
	// 	secretPath:   "",
	// 	hostFilePath: "",
	// 	manager:      p,
	// 	fileSystem:   f,
	// }
	//
	// for _, tt := range tests {
	// 	for _, value := range tt.dockerClientMockArgs {
	// 		a.On(value.method, value.args...).Return(value.paramReturned...)
	// 	}
	// 	for _, value := range tt.fileSystemMockArgs {
	// 		f.On(value.method, value.args...).Return(value.paramReturned...)
	// 	}
	// 	for _, value := range tt.proxyManagerMockArgs {
	// 		p.On(value.method, value.args...).Return(value.paramReturned...)
	// 	}
	// 	err := d.ApplyServiceRoutes(tt.args.ctx, tt.args.projectID, tt.args.serviceID, tt.args.routes)
	// 	assertObj.Equal(tt.wantErr, err != nil, fmt.Sprintf("Test name (%s) Got error (%v)", tt.name, err))
	// }
	// a.AssertExpectations(t)
	// f.AssertExpectations(t)
	// p.AssertExpectations(t)
}

func TestDocker_DeleteService(t *testing.T) {

	type mockArg struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	type args struct {
		ctx       context.Context
		projectID string
		serviceID string
		version   string
	}
	tests := []struct {
		name                 string
		dockerClientMockArgs []mockArg
		fileSystemMockArgs   []mockArg
		proxyManagerMockArgs []mockArg
		args                 args
		wantErr              bool
	}{
		{
			name: "service deleted successfully -- service with specific version deleted",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject1--dummyService1--dummyVersion1")), All: true}},
					paramReturned: []interface{}{[]types.Container{{ID: "dummyService1-container"}}, nil},
				},
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject1--dummyService1")), All: true}},
					paramReturned: []interface{}{[]types.Container{}, nil},
				},
				{
					method:        "ContainerRemove",
					args:          []interface{}{context.Background(), "dummyService1-container", types.ContainerRemoveOptions{Force: true}},
					paramReturned: []interface{}{nil},
				},
			},
			fileSystemMockArgs: []mockArg{
				{
					method:        "NewHostFile",
					args:          []interface{}{},
					paramReturned: []interface{}{&txeh.Hosts{}, nil},
				},
				{
					method:        "RemoveTempSecretsFolder",
					args:          []interface{}{"dummyProject1", "dummyService1", "dummyVersion1"},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetInternalServiceDomain("dummyProject1", "dummyService1", "dummyVersion1")},
					paramReturned: []interface{}{},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetServiceDomain("dummyProject1", "dummyService1")},
					paramReturned: []interface{}{},
				},
				{
					method:        "SaveHostFile",
					args:          []interface{}{&txeh.Hosts{}},
					paramReturned: []interface{}{nil},
				},
			},
			proxyManagerMockArgs: []mockArg{
				{
					method:        "DeleteServiceRoutes",
					args:          []interface{}{"dummyProject1", "dummyService1"},
					paramReturned: []interface{}{nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject1",
				serviceID: "dummyService1",
				version:   "dummyVersion1",
			},
		},
		{
			name: "service deleted successfully -- service with only project id -- delete all container having similar project id",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "dummyProject1")), All: true}},
					paramReturned: []interface{}{[]types.Container{{ID: "dummyProject1-container1"}, {ID: "dummyProject1-container2"}}, nil},
				},
				{ // used by check if last service
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject1--")), All: true}},
					paramReturned: []interface{}{[]types.Container{}, nil},
				},
				{
					method:        "ContainerRemove",
					args:          []interface{}{context.Background(), "dummyProject1-container1", types.ContainerRemoveOptions{Force: true}},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "ContainerRemove",
					args:          []interface{}{context.Background(), "dummyProject1-container2", types.ContainerRemoveOptions{Force: true}},
					paramReturned: []interface{}{nil},
				},
			},
			fileSystemMockArgs: []mockArg{
				{
					method:        "NewHostFile",
					args:          []interface{}{},
					paramReturned: []interface{}{&txeh.Hosts{}, nil},
				},
				{
					method:        "RemoveTempSecretsFolder",
					args:          []interface{}{"dummyProject1", "", ""},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetInternalServiceDomain("dummyProject1", "", "")},
					paramReturned: []interface{}{},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetServiceDomain("dummyProject1", "")},
					paramReturned: []interface{}{},
				},
				{
					method:        "SaveHostFile",
					args:          []interface{}{&txeh.Hosts{}},
					paramReturned: []interface{}{nil},
				},
			},
			proxyManagerMockArgs: []mockArg{
				{
					method:        "DeleteServiceRoutes",
					args:          []interface{}{"dummyProject1", ""},
					paramReturned: []interface{}{nil},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject1",
			},
		},
		{
			name: "service deletion failed -- unable to list containers",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject2--dummyService2--dummyVersion2")), All: true}},
					paramReturned: []interface{}{nil, fmt.Errorf("no container found with applied filter opeartion")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject2",
				serviceID: "dummyService2",
				version:   "dummyVersion2",
			},
			wantErr: true,
		},
		{
			name: "service deletion failed -- unable to remove specified container",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject3--dummyService3--dummyVersion3")), All: true}},
					paramReturned: []interface{}{[]types.Container{{ID: "dummyService3-container"}}, nil},
				},
				{
					method:        "ContainerRemove",
					args:          []interface{}{context.Background(), "dummyService3-container", types.ContainerRemoveOptions{Force: true}},
					paramReturned: []interface{}{fmt.Errorf("unable to remove specified conatainer")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject3",
				serviceID: "dummyService3",
				version:   "dummyVersion3",
			},
			wantErr: true,
		},
		// {
		// 	name: "service deletion failed -- unable to create new host file",
		// 	dockerClientMockArgs: []mockArg{
		// 		{
		// 			method:        "ContainerList",
		// 			args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject4--dummyService4--dummyVersion4")), All: true}},
		// 			paramReturned: []interface{}{[]types.Container{{ID: "dummyService4-container"}}, nil},
		// 		},
		// 		{
		// 			method:        "ContainerRemove",
		// 			args:          []interface{}{context.Background(), "dummyService4-container", types.ContainerRemoveOptions{Force: true}},
		// 			paramReturned: []interface{}{nil},
		// 		},
		// 	},
		// 	fileSystemMockArgs: []mockArg{
		// 		{
		// 			method:        "NewHostFile",
		// 			args:          []interface{}{},
		// 			paramReturned: []interface{}{nil, fmt.Errorf("unable to create host file")},
		// 		},
		// 	},
		// 	args: args{
		// 		ctx:       context.Background(),
		// 		projectID: "dummyProject4",
		// 		serviceID: "dummyService4",
		// 		version:   "dummyVersion4",
		// 	},
		// 	wantErr: true,
		// },
		{
			name: "service deletion failed -- unable to check if last service",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject4--dummyService4--dummyVersion4")), All: true}},
					paramReturned: []interface{}{[]types.Container{{ID: "dummyService4-container"}}, nil},
				},
				{
					method:        "ContainerRemove",
					args:          []interface{}{context.Background(), "dummyService4-container", types.ContainerRemoveOptions{Force: true}},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject4--dummyService4")), All: true}},
					paramReturned: []interface{}{nil, fmt.Errorf("unable to list specified container with fileter options")},
				},
			},
			fileSystemMockArgs: []mockArg{
				{
					method:        "NewHostFile",
					args:          []interface{}{},
					paramReturned: []interface{}{&txeh.Hosts{}, nil},
				},
				{
					method:        "RemoveTempSecretsFolder",
					args:          []interface{}{"dummyProject4", "dummyService4", "dummyVersion4"},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetInternalServiceDomain("dummyProject4", "dummyService4", "dummyVersion4")},
					paramReturned: []interface{}{},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject4",
				serviceID: "dummyService4",
				version:   "dummyVersion4",
			},
			wantErr: true,
		},
		{
			name: "service deletion failed -- proxy manager unable to delete service routes",
			dockerClientMockArgs: []mockArg{
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject5--dummyService5--dummyVersion5")), All: true}},
					paramReturned: []interface{}{[]types.Container{{ID: "dummyService5-container"}}, nil},
				},
				{
					method:        "ContainerRemove",
					args:          []interface{}{context.Background(), "dummyService5-container", types.ContainerRemoveOptions{Force: true}},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "ContainerList",
					args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud-dummyProject5--dummyService5")), All: true}},
					paramReturned: []interface{}{[]types.Container{}, nil},
				},
			},
			fileSystemMockArgs: []mockArg{
				{
					method:        "NewHostFile",
					args:          []interface{}{},
					paramReturned: []interface{}{&txeh.Hosts{}, nil},
				},
				{
					method:        "RemoveTempSecretsFolder",
					args:          []interface{}{"dummyProject5", "dummyService5", "dummyVersion5"},
					paramReturned: []interface{}{nil},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetInternalServiceDomain("dummyProject5", "dummyService5", "dummyVersion5")},
					paramReturned: []interface{}{},
				},
				{
					method:        "RemoveHostFromHostFile",
					args:          []interface{}{&txeh.Hosts{}, utils.GetServiceDomain("dummyProject5", "dummyService5")},
					paramReturned: []interface{}{},
				},
			},
			proxyManagerMockArgs: []mockArg{
				{
					method:        "DeleteServiceRoutes",
					args:          []interface{}{"dummyProject5", "dummyService5"},
					paramReturned: []interface{}{fmt.Errorf("unable to delete service routes of specified service")},
				},
			},
			args: args{
				ctx:       context.Background(),
				projectID: "dummyProject5",
				serviceID: "dummyService5",
				version:   "dummyVersion5",
			},
			wantErr: true,
		},
	}

	t.Parallel()
	assertObj := assert.New(t)
	a := &dockerClientMock{}
	f := &fileSystemMock{}
	p := &proxyMangerMock{}
	d := &Docker{
		client:       a,
		auth:         nil,
		artifactAddr: "",
		secretPath:   "",
		hostFilePath: "",
		manager:      p,
		fileSystem:   f,
	}

	for _, tt := range tests {
		for _, value := range tt.dockerClientMockArgs {
			a.On(value.method, value.args...).Return(value.paramReturned...)
		}
		for _, value := range tt.fileSystemMockArgs {
			f.On(value.method, value.args...).Return(value.paramReturned...)
		}
		for _, value := range tt.proxyManagerMockArgs {
			p.On(value.method, value.args...).Return(value.paramReturned...)
		}
		err := d.DeleteService(tt.args.ctx, tt.args.projectID, tt.args.serviceID, tt.args.version)
		assertObj.Equal(tt.wantErr, err != nil, fmt.Sprintf("Test name (%s) Got error (%v)", tt.name, err))
	}
	a.AssertExpectations(t)
	f.AssertExpectations(t)
	p.AssertExpectations(t)
}

func TestDocker_GetServiceRoutes(t *testing.T) {
	type fields struct {
		client       dockerClient
		auth         *auth.Module
		artifactAddr string
		secretPath   string
		hostFilePath string
		manager      proxyManager
		fileSystem   fileSystem
	}
	type args struct {
		in0       context.Context
		projectID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]model.Routes
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Docker{
				client:       tt.fields.client,
				auth:         tt.fields.auth,
				artifactAddr: tt.fields.artifactAddr,
				secretPath:   tt.fields.secretPath,
				hostFilePath: tt.fields.hostFilePath,
				manager:      tt.fields.manager,
				fileSystem:   tt.fields.fileSystem,
			}
			got, err := d.GetServiceRoutes(tt.args.in0, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServiceRoutes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocker_GetServices(t *testing.T) {
	type fields struct {
		client       dockerClient
		auth         *auth.Module
		artifactAddr string
		secretPath   string
		hostFilePath string
		manager      proxyManager
		fileSystem   fileSystem
	}
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Docker{
				client:       tt.fields.client,
				auth:         tt.fields.auth,
				artifactAddr: tt.fields.artifactAddr,
				secretPath:   tt.fields.secretPath,
				hostFilePath: tt.fields.hostFilePath,
				manager:      tt.fields.manager,
				fileSystem:   tt.fields.fileSystem,
			}
			got, err := d.GetServices(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocker_checkIfLastService(t *testing.T) {
	type fields struct {
		client       dockerClient
		auth         *auth.Module
		artifactAddr string
		secretPath   string
		hostFilePath string
		manager      proxyManager
		fileSystem   fileSystem
	}
	type args struct {
		ctx       context.Context
		projectID string
		serviceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Docker{
				client:       tt.fields.client,
				auth:         tt.fields.auth,
				artifactAddr: tt.fields.artifactAddr,
				secretPath:   tt.fields.secretPath,
				hostFilePath: tt.fields.hostFilePath,
				manager:      tt.fields.manager,
				fileSystem:   tt.fields.fileSystem,
			}
			got, err := d.checkIfLastService(tt.args.ctx, tt.args.projectID, tt.args.serviceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfLastService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfLastService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocker_createContainer(t *testing.T) {
	type fields struct {
		client       dockerClient
		auth         *auth.Module
		artifactAddr string
		secretPath   string
		hostFilePath string
		manager      proxyManager
		fileSystem   fileSystem
	}
	type args struct {
		ctx           context.Context
		index         int
		task          model.Task
		service       *model.Service
		overridePorts []model.Port
		cName         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Docker{
				client:       tt.fields.client,
				auth:         tt.fields.auth,
				artifactAddr: tt.fields.artifactAddr,
				secretPath:   tt.fields.secretPath,
				hostFilePath: tt.fields.hostFilePath,
				manager:      tt.fields.manager,
				fileSystem:   tt.fields.fileSystem,
			}
			got, got1, err := d.createContainer(tt.args.ctx, tt.args.index, tt.args.task, tt.args.service, tt.args.overridePorts, tt.args.cName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createContainer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("createContainer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewDockerDriver(t *testing.T) {
	type args struct {
		auth         *auth.Module
		artifactAddr string
	}
	tests := []struct {
		name    string
		args    args
		want    *Docker
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDockerDriver(tt.args.auth, tt.args.artifactAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDockerDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDockerDriver() got = %v, want %v", got, tt.want)
			}
		})
	}
}
