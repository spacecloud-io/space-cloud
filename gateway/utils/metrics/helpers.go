package metrics

import (
	"fmt"
	"strings"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

const (
	eventingModule      = "eventing"
	fileModule          = "file"
	databaseModule      = "db"
	remoteServiceModule = "remote-service" // aka remote service
	notApplicable       = "na"
)

func newMetrics() *metrics {
	return &metrics{}
}

func getModuleName(key string) string {
	return strings.Split(key, ":")[0]
}

func generateEventingKey(project, eventingType string) string {
	return fmt.Sprintf("%s:%s:%s", eventingModule, project, eventingType)
}

func parseEventingKey(key string) (module, project, eventingType string) {
	v := strings.Split(key, ":")
	return v[0], v[1], v[2]
}

func generateFunctionKey(project, serviceName, functionName string) string {
	return fmt.Sprintf("%s:%s:%s:%s", remoteServiceModule, project, serviceName, functionName)
}

func parseFunctionKey(key string) (module, project, remoteServiceName, endpointName string) {
	v := strings.Split(key, ":")
	return v[0], v[1], v[2], v[3]
}

func generateDatabaseKey(project, dbAlias, tableName string) string {
	return fmt.Sprintf("%s:%s:%s:%s", databaseModule, project, dbAlias, tableName)
}

func parseDatabaseKey(key string) (module, project, dbAlias, tableName string) {
	v := strings.Split(key, ":")
	return v[0], v[1], v[2], v[3]
}

func generateFileKey(project, storeType string) string {
	return fmt.Sprintf("%s:%s:%s", fileModule, project, storeType)
}

func parseFileKey(key string) (module, project, fileStoreType string) {
	v := strings.Split(key, ":")
	return v[0], v[1], v[2]
}

func (m *Module) createFileDocuments(key string, metrics *metricOperations, t string) []interface{} {
	docs := make([]interface{}, 0)
	module, projectName, storeType := parseFileKey(key)
	if metrics.create > 0 {
		docs = append(docs, m.createDocument(projectName, storeType, notApplicable, module, utils.Create, metrics.create, t))
	}

	if metrics.read > 0 {
		docs = append(docs, m.createDocument(projectName, storeType, notApplicable, module, utils.Read, metrics.read, t))
	}

	if metrics.delete > 0 {
		docs = append(docs, m.createDocument(projectName, storeType, notApplicable, module, utils.Delete, metrics.delete, t))
	}

	if metrics.list > 0 {
		docs = append(docs, m.createDocument(projectName, storeType, notApplicable, module, utils.List, metrics.list, t))
	}

	return docs
}

func (m *Module) createCrudDocuments(key string, value *metricOperations, t string) []interface{} {
	docs := make([]interface{}, 0)
	module, projectName, dbAlias, tableName := parseDatabaseKey(key)
	if value.create > 0 {
		docs = append(docs, m.createDocument(projectName, dbAlias, tableName, module, utils.Create, value.create, t))
	}

	if value.read > 0 {
		docs = append(docs, m.createDocument(projectName, dbAlias, tableName, module, utils.Read, value.read, t))
	}

	if value.update > 0 {
		docs = append(docs, m.createDocument(projectName, dbAlias, tableName, module, utils.Update, value.update, t))
	}

	if value.delete > 0 {
		docs = append(docs, m.createDocument(projectName, dbAlias, tableName, module, utils.Delete, value.delete, t))
	}

	return docs
}

func (m *Module) createEventDocument(key string, count uint64, t string) []interface{} {
	module, projectName, eventingType := parseEventingKey(key)
	docs := make([]interface{}, 0)
	if count > 0 {
		docs = append(docs, m.createDocument(projectName, notApplicable, notApplicable, module, utils.OperationType(eventingType), count, t))
	}
	return docs
}

func (m *Module) createFunctionDocument(key string, count uint64, t string) []interface{} {
	module, projectName, serviceName, functionName := parseFunctionKey(key)
	docs := make([]interface{}, 0)
	if count > 0 {
		docs = append(docs, m.createDocument(projectName, serviceName, functionName, module, "calls", count, t))
	}
	return docs
}

func (m *Module) createDocument(project, driver, subType, module string, op utils.OperationType, count uint64, t string) interface{} {
	return map[string]interface{}{
		"id":         ksuid.New().String(),
		"project_id": project,
		"module":     module,
		"type":       op,
		"sub_type":   subType,
		"ts":         t,
		"count":      count,
		"driver":     driver,
		"node_id":    m.nodeID,
		"cluster_id": m.clusterID,
	}
}
