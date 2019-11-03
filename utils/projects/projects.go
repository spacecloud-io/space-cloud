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
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/modules/userman"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/graphql"
	"github.com/spaceuptech/space-cloud/utils/metrics"
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
	Eventing       *eventing.Module
	Graph          *graphql.Module
}

// Projects is the stub to manage the state of the various modules
type Projects struct {
	lock               sync.RWMutex
	nodeID             string
	removeProjectScope bool
	projects           map[string]*ProjectState
	h                  *driver.Handler

	// Global managers
	syncMan  *syncman.Manager
	adminMan *admin.Manager
	metrics  *metrics.Module
}

// New creates a new Projects instance
func New(nodeID string, removeProjectScope bool, h *driver.Handler,
	adminMan *admin.Manager, syncMan *syncman.Manager, metrics *metrics.Module) *Projects {
	return &Projects{nodeID: nodeID, removeProjectScope: removeProjectScope, projects: map[string]*ProjectState{}, h: h,
		syncMan: syncMan, adminMan: adminMan, metrics: metrics}
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

	// Create the fundamental modules
	c := crud.Init(p.h)

	fn := functions.Init(p.syncMan)

	a := auth.Init(c, fn, p.removeProjectScope)
	a.SetMakeHttpRequest(p.syncMan.MakeHTTPRequest)

	// Initialise the eventing module and set the crud module hooks
	e := eventing.New(a, c, fn, p.adminMan, p.syncMan)

	c.SetHooks(&model.CrudHooks{
		Create: e.HandleCreateIntent,
		Update: e.HandleUpdateIntent,
		Delete: e.HandleDeleteIntent,
		Batch:  e.HandleBatchIntent,
		Stage:  e.HandleStage,
	}, p.metrics.AddDBOperation)

	rt, err := realtime.Init(p.nodeID, e, a, c, p.metrics, p.syncMan)
	if err != nil {
		return nil, err
	}

	u := userman.Init(c, a)
	f := filestore.Init(a)
	graph := graphql.New(a, c, fn)

	state := &ProjectState{Crud: c, Functions: fn, Auth: a, UserManagement: u, FileStore: f, Realtime: rt,
		Eventing: e, Graph: graph}

	p.projects[project] = state

	return state, nil
}

func (p *Projects) GetProjectIDs() []string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var ids []string
	for id := range p.projects {
		ids = append(ids, id)
	}

	return ids
}
