package file

import (
	"io/ioutil"
	"path/filepath"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (l *File) loadConfiguration() (map[string][]*unstructured.Unstructured, map[string][]*unstructured.Unstructured, error) {
	files, err := ioutil.ReadDir(l.Path)
	if err != nil {
		l.logger.Error("Unable to read config files from directory", zap.String("dir", l.Path), zap.Error(err))
		return nil, nil, err
	}

	configuration := map[string][]*unstructured.Unstructured{}
	configurationN := map[string][]*unstructured.Unstructured{}

	for _, file := range files {
		arr, err := utils.ReadSpecObjectsFromFile(filepath.Join(l.Path, file.Name()))
		if err != nil {
			l.logger.Error("Unable to parse config file", zap.String("file", file.Name()), zap.Error(err))
			return nil, nil, err
		}

		for _, spec := range arr {
			kind := spec.GetKind()
			key := source.GetModuleName(spec.GetAPIVersion(), spec.GetKind())
			configuration[kind] = append(configuration[kind], spec)
			configurationN[key] = append(configurationN[key], spec)
		}
	}

	return configuration, configurationN, nil
}
