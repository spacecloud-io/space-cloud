package static

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

const (
	defaultDirPath   = "./public"
	defaultUrlPrefix = "/"
)

// Module is responsible for Static
type Module struct {
	sync.RWMutex
	Enabled   bool
	Path      string
	URLPrefix string
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
	return &Module{Path: defaultDirPath, URLPrefix: defaultUrlPrefix, Gzip:false}
}

// SetConfig set the config required by the Static module
func (m *Module) SetConfig(s *config.Static) error {
	m.Lock()
	defer m.Unlock()

	if s == nil || !s.Enabled {
		m.Enabled = false
		return nil
	}

	m.Gzip = s.Gzip

	m.Path = s.Path
	if m.Path == "" {
		m.Path = defaultDirPath
	}

	m.URLPrefix = s.URLPrefix
	if m.URLPrefix == "" {
		m.URLPrefix = defaultUrlPrefix
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
