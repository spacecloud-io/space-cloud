package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (f *File) loadConfiguration() error {
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		f.logger.Error("Unable to read config files from directory", zap.String("dir", f.path), zap.Error(err))
		return err
	}

	configuration := common.ConfigType{}

	for _, file := range files {
		arr, err := utils.ReadSpecObjectsFromFile(filepath.Join(f.path, file.Name()))
		if err != nil {
			f.logger.Error("Unable to parse config file", zap.String("file", file.Name()), zap.Error(err))
			return err
		}

		for _, spec := range arr {
			gvr := schema.GroupVersionResource{
				Group:    spec.GroupVersionKind().Group,
				Version:  spec.GroupVersionKind().Version,
				Resource: utils.Pluralize(spec.GetKind())}
			key := source.GetModuleName(gvr)
			configuration[key] = append(configuration[key], spec)
		}
	}

	f.setConfig(configuration)
	return nil
}

// reorganizeFileStructure organizes the directory where configurations are stored.
// Directory is organized in such a way that all the configurations of a source are
// stored in a single source file.
func (f *File) reorganizeFileStructure() error {
	// Get source file names which is same as "Resource" of GVR of source.
	registeredSources := source.GetRegisteredSources()
	sourceFileNames := make(map[string]struct{})
	for _, src := range registeredSources {
		fileName := generateYAMLFileName(src.Group, src.Resource)
		sourceFileNames[fileName] = struct{}{}
	}

	// Get all files in the config directory.
	files, err := ioutil.ReadDir(f.path)
	if err != nil {
		f.logger.Error("Unable to read config files from directory", zap.String("dir", f.path), zap.Error(err))
		return err
	}

	// Check if file is a valid source file. If not, then delete it and append
	// its specs to the valid source file.
	for _, file := range files {
		if _, ok := sourceFileNames[file.Name()]; !ok {
			// Parse the config file
			arr, err := utils.ReadSpecObjectsFromFile(filepath.Join(f.path, file.Name()))
			if err != nil {
				f.logger.Error("Unable to parse config file", zap.String("file", file.Name()), zap.Error(err))
				return err
			}

			for _, spec := range arr {
				// Get spec in bytes
				data, err := utils.GetBytesFromSpec(spec)
				if err != nil {
					return err
				}

				// source file
				fileName := generateYAMLFileName(spec.GroupVersionKind().Group, utils.Pluralize(spec.GetKind()))

				// Check if the source file exists
				_, err = os.Stat(filepath.Join(f.path, fileName))
				// If source file does not exists create a new one and
				// write this spec into this file
				if os.IsNotExist(err) {
					err := os.WriteFile(filepath.Join(f.path, fileName), data, 0777)
					if err != nil {
						return err
					}
				} else {
					// If source file exists append the spec.
					err := utils.AppendToFile(filepath.Join(f.path, fileName), data)
					if err != nil {
						return err
					}
				}
			}

			// delete the file
			if err = os.Remove(filepath.Join(f.path, file.Name())); err != nil {
				return fmt.Errorf("could not delete file: %v", err)
			}
		}
	}
	return nil
}

func (f *File) setConfig(newConfig common.ConfigType) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.configuration = newConfig
}

func (f *File) getConfig() common.ConfigType {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.configuration
}

func generateYAMLFileName(group, resource string) string {
	fileName := group + "." + resource + ".yaml"
	return fileName
}

func (f *File) persistConfig(gvr schema.GroupVersionResource, data []byte) error {
	fileName := generateYAMLFileName(gvr.Group, gvr.Resource)
	// Check if the source file exists
	_, err := os.Stat(filepath.Join(f.path, fileName))
	// If source file does not exists create a new one and
	// write this spec into this file
	if os.IsNotExist(err) {
		err := os.WriteFile(filepath.Join(f.path, fileName), data, 0777)
		if err != nil {
			return err
		}
	} else {
		// If source file exists append the spec.
		err := utils.AppendToFile(filepath.Join(f.path, fileName), data)
		if err != nil {
			return err
		}
	}
	return nil
}
