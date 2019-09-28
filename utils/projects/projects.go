package projects

import (
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/modules/eventing"
	"github.com/spaceuptech/space-cloud/modules/filestore"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/modules/pubsub"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/userman"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/graphql"
	"github.com/spaceuptech/space-cloud/utils/syncman"
)

// ProjectState holds the module state of a project
type ProjectState struct {
	//Config         *config.Project
	Auth           *auth.Module
	Crud           *crud.Module
	UserManagement *userman.Module
	FileStore      *filestore.Module
	Functions      *functions.Module
	Realtime       *realtime.Module
	Pubsub         *pubsub.Module
	Eventing       *eventing.Module
	Graph          *graphql.Module
}

// Projects is the stub to manage the state of the various modules
type Projects struct {
	lock     sync.RWMutex
	nodeID   string
	projects map[string]*ProjectState
	h        *driver.Handler

	// Global managers
	syncMan  *syncman.SyncManager
	adminMan *admin.Manager
}

// New creates a new Projects instance
func New(nodeID string, h *driver.Handler, adminMan *admin.Manager, syncMan *syncman.SyncManager) *Projects {
	return &Projects{nodeID: nodeID, projects: map[string]*ProjectState{}, h: h, syncMan: syncMan, adminMan: adminMan}
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

// Iter iterates over all the projects and passes it in the provided function.
// Iteration stops if the function returns false
func (p *Projects) Iter(fn func(string, *ProjectState) bool) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for project, state := range p.projects {
		if !fn(project, state) {
			return false
		}
	}

	return true
}

// NewProject creates a new project with all modules in the default state.
// It will overwrite the existing project if any
func (p *Projects) NewProject(project string) (*ProjectState, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	c := crud.Init(p.h)
	f, err := functions.Init()
	if err != nil {
		return nil, err
	}

	a := auth.Init(c, f)
	u := userman.Init(c, a)
	file := filestore.Init(a)
	eventing := eventing.New(c, f, p.syncMan)

	c.SetHooks(&model.CrudHooks{
		Create: eventing.HandleCreateIntent,
		Update: eventing.HandleUpdateIntent,
		Delete: eventing.HandleDeleteIntent,
		Batch:  eventing.HandleBatchIntent,
		Stage:  eventing.HandleStage,
	})

	r, err := realtime.Init(p.nodeID, eventing, a, c, f, p.adminMan, p.syncMan)
	if err != nil {
		return nil, err
	}

	graph := graphql.New(a, c, f)

	state := &ProjectState{Crud: c, Functions: f, Auth: a, UserManagement: u, FileStore: file, Realtime: r,
		Eventing: eventing, Graph: graph}
	p.projects[project] = state

	return state, nil
}
