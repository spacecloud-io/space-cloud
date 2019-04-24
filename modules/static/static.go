package static

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

const (
	defaultDirPath   = "./static"
	defaultUrlPrefix = "/static/"
)

// Module is responsible for Static
type Module struct {
	sync.RWMutex
	Enabled   bool
	Path      string
	UrlPrefix string
	Gzip      bool
}

type SpaFileSystem struct {
	root http.FileSystem
}

func (fs *SpaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		log.Panicln("[ERROR: Static module]: ", err)
	}
	return f, err
}

// Init returns a new instance of the Static module wit default values
func Init() *Module {
	return &Module{Path: defaultDirPath, UrlPrefix: defaultUrlPrefix, Gzip:false}
}

// SetConfig set the config required by the Static module
func (m *Module) SetConfig(s *config.Static) error {
	m.Lock()
	defer m.Unlock()

	if s == nil || !s.Enabled {
		m.Enabled = false
		return nil
	}

	if s.Gzip {
		m.Gzip = true
	}

	if s.Path != "" {
		m.Path = s.Path
	}

	if s.UrlPrefix != "" {
		m.Path = s.Path
	}

	m.Enabled = true
	return nil
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.Enabled
}

func (m *Module) getDirPath() string {
	m.RLock()
	defer m.RUnlock()

	return m.Path
}
