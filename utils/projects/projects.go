package projects

import (
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/filestore"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/userman"
)

// ProjectState holds the module state of a project
type ProjectState struct {
	Auth           *auth.Module
	Crud           *crud.Module
	UserManagement *userman.Module
	FileStore      *filestore.Module
	//Config         *config.Project
	Functions *functions.Module
	Realtime  *realtime.Module
}

// Projects is the stub to manage the state of the various modules
type Projects struct {
	lock     sync.RWMutex
	projects map[string]*ProjectState
}

// New creates a new Projects instance
func New() *Projects {
	return &Projects{projects: map[string]*ProjectState{}}
}

// LoadProject returns the state of the project specified
func (p *Projects) LoadProject(project string) (*ProjectState, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if state, p := p.projects[project]; p {
		return state, nil
	}

	return nil, errors.New("project not found in server state")
}

// DeleteProject deletes a single project
func (p *Projects) DeleteProject(project string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.projects, project)
}

// NewProject creates a new project with all modules in the default state.
// It will overwrite the existing project if any
func (p *Projects) NewProject(project string) *ProjectState {
	p.lock.Lock()
	defer p.lock.Unlock()

	c := crud.Init()
	f := functions.Init()
	a := auth.Init(c, nil)
	u := userman.Init(c, a)
	file := filestore.Init()
	r := realtime.Init()

	state := &ProjectState{Crud: c, Functions: f, Auth: a, UserManagement: u, FileStore: file, Realtime: r}
	p.projects[project] = state

	return state
}
