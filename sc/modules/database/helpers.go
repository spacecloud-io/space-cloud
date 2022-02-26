package database

import (
	"fmt"
	"strings"

	"github.com/spacecloud-io/space-cloud/config"
)

// CombineDBConfigKey joing project id and db alias to form the db config key
func CombineDBConfigKey(projectID, dbAlias string) string {
	return fmt.Sprintf("%s---%s", projectID, dbAlias)
}

// SplitDBConfigKey split the db config key into project id and dbAlias
func SplitDBConfigKey(key string) (project, dbAlias string) {
	arr := strings.Split(key, "---")
	return arr[0], arr[1]
}

func generateUniqueDBKey(projectID string, c *config.DatabaseConfig) string {
	return fmt.Sprintf("%s---%s--%v", CombineDBConfigKey(projectID, c.DbAlias), c.DBName, c.DriverConf)
}
