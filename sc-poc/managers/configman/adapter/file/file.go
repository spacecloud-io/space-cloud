package file

import (
	"context"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spacecloud-io/space-cloud/managers/configman/adapter"
	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"go.uber.org/zap"
)

type File struct {
	Path          string
	logger        *zap.Logger
	lock          sync.RWMutex
	configuration common.ConfigType
}

// MakeFileAdapter returns the File adapter object.
func MakeFileAdapter(path string) adapter.Adapter {
	logger, _ := zap.NewDevelopment()
	file := &File{
		Path:          path,
		logger:        logger,
		configuration: make(common.ConfigType),
	}
	return file
}

// GetRawConfig returns the final config.
func (f *File) GetRawConfig() (common.ConfigType, error) {
	// Load SC config file from file system
	if err := f.loadConfiguration(); err != nil {
		return nil, err
	}

	return f.getConfig(), nil
}

// Run watches the files indefinitely.
func (file *File) Run(ctx context.Context) (chan common.ConfigType, error) {
	cfgChan := make(chan common.ConfigType)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return cfgChan, err
	}

	go file.watchEvents(watcher, cfgChan)

	err = watcher.Add(file.Path)
	if err != nil {
		return cfgChan, err
	}

	return cfgChan, nil
}

func (f *File) watchEvents(watcher *fsnotify.Watcher, cfgChan chan common.ConfigType) {
	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				f.logger.Error("channel to watching config file closed")
				break
			}

			resp, err := f.GetRawConfig()
			if err != nil {
				f.logger.Error("reloading config", zap.Error(err))
				break
			}

			cfgChan <- resp
		case err := <-watcher.Errors:
			f.logger.Error("issue with file watcher", zap.Error(err))
		}
	}
}

// Interface guard
var (
	_ adapter.Adapter = (*File)(nil)
)
