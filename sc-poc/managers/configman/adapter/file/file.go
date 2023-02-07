package file

import (
	"context"
	"encoding/json"

	"github.com/fsnotify/fsnotify"
	"github.com/spacecloud-io/space-cloud/managers/configman/adapter"
	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"go.uber.org/zap"
)

type File struct {
	Path   string
	logger *zap.Logger
}

// MakeFileAdapter returns the File adapter object.
func MakeFileAdapter(path string) adapter.Adapter {
	logger, _ := zap.NewDevelopment()
	file := &File{
		Path:   path,
		logger: logger,
	}
	return file
}

// GetRawConfig returns the final caddy config in bytes.
func (f *File) GetRawConfig() ([]byte, error) {
	// Load SC config file from file system
	configuration, err := f.loadConfiguration()
	if err != nil {
		return nil, err
	}

	// Load the new caddy config
	config, err := common.PrepareConfig(configuration)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(config, "", "  ")
}

// Run watches the files indefinitely.
func (file *File) Run(ctx context.Context) (chan []byte, error) {
	cfgChan := make(chan []byte)
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

func (f *File) watchEvents(watcher *fsnotify.Watcher, cfgChan chan []byte) {
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
			break
		}
	}
}

// Interface guard
var (
	_ adapter.Adapter = (*File)(nil)
)
