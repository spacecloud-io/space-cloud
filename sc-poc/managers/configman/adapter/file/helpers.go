package file

import (
	"io/ioutil"
	"path/filepath"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (l *File) loadConfiguration() (common.ConfigType, error) {
	files, err := ioutil.ReadDir(l.Path)
	if err != nil {
		l.logger.Error("Unable to read config files from directory", zap.String("dir", l.Path), zap.Error(err))
		return nil, err
	}

	configuration := common.ConfigType{}

	for _, file := range files {
		arr, err := utils.ReadSpecObjectsFromFile(filepath.Join(l.Path, file.Name()))
		if err != nil {
			l.logger.Error("Unable to parse config file", zap.String("file", file.Name()), zap.Error(err))
			return nil, err
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

	return configuration, nil
}
