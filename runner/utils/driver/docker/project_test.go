package docker

import (
	"context"
	"fmt"
	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocker_CreateProject(t *testing.T) {

	type mockArg struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	type args struct {
		ctx     context.Context
		project *model.Project
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
			name: "successfully deleted project",
			fileSystemMockArgs: []mockArg{
				{
					method:        "CreateProjectDirectory",
					args:          []interface{}{"dummyProject1"},
					paramReturned: []interface{}{nil},
				},
			},
			args: args{
				ctx: context.Background(),
				project: &model.Project{
					ID:          "dummyProject1",
					Environment: "dummyEnvironment1",
				},
			},
			wantErr: false,
		},
		{
			name: "successfully deleted project -- unable to delete project directory",
			fileSystemMockArgs: []mockArg{
				{
					method:        "CreateProjectDirectory",
					args:          []interface{}{"dummyProject2"},
					paramReturned: []interface{}{fmt.Errorf("unable to delete specified project directory")},
				},
			},
			args: args{
				ctx: context.Background(),
				project: &model.Project{
					ID:          "dummyProject2",
					Environment: "dummyEnvironment2",
				},
			},
			wantErr: true,
		},
	}

	t.Parallel()
	assertObj := assert.New(t)
	f := &fileSystemMock{}
	d := &Docker{
		fileSystem: f,
	}

	for _, tt := range tests {
		for _, value := range tt.fileSystemMockArgs {
			f.On(value.method, value.args...).Return(value.paramReturned...)
		}
		err := d.CreateProject(tt.args.ctx, tt.args.project)
		assertObj.Equal(tt.wantErr, err != nil, fmt.Sprintf("Test name (%s) Got error (%v)", tt.name, err))
	}
	f.AssertExpectations(t)
}

func TestDocker_DeleteProject(t *testing.T) {
	//
	// type mockArg struct {
	// 	method        string
	// 	args          []interface{}
	// 	paramReturned []interface{}
	// }
	// type args struct {
	// 	ctx       context.Context
	// 	projectID string
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
	// 		name: "successfully deleted project",
	// 		dockerClientMockArgs: []mockArg{
	// 			{
	// 				method:        "DeleteService",
	// 				args:          []interface{}{context.Background(), "dummyProject1", "", ""},
	// 				paramReturned: []interface{}{nil},
	// 			},
	// 			{
	// 				method:        "ContainerList",
	// 				args:          []interface{}{context.Background(), types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "dummyProject1")), All: true}},
	// 				paramReturned: []interface{}{[]types.Container{{ID: "dummyProject1-container1"}, {ID: "dummyProject1-container2"}}, nil},
	// 			},
	// 		},
	// 		fileSystemMockArgs: []mockArg{
	// 			{
	// 				method:        "RemoveProjectDirectory",
	// 				args:          []interface{}{"dummyProject1"},
	// 				paramReturned: []interface{}{nil},
	// 			},
	// 		},
	// 		args: args{
	// 			ctx:       context.Background(),
	// 			projectID: "dummyProject1",
	// 		},
	// 		wantErr: false,
	// 	},
	// }
	// t.Parallel()
	// assertObj := assert.New(t)
	// a := &dockerClientMock{}
	// f := &fileSystemMock{}
	// d := &Docker{
	// 	client:       a,
	// 	auth:         nil,
	// 	artifactAddr: "",
	// 	secretPath:   "",
	// 	hostFilePath: "",
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
	// 	err := d.DeleteProject(tt.args.ctx, tt.args.projectID)
	// 	assertObj.Equal(tt.wantErr, err != nil, fmt.Sprintf("Test name (%s) Got error (%v)", tt.name, err))
	// }
	// a.AssertExpectations(t)
	// f.AssertExpectations(t)
}
