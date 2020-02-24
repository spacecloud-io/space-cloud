package crud

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (m *Module) createBatch(project, dbAlias, col string, doc interface{}) (int64, error) {
	response := make(batchResponseChan, 1)
	defer close(response)

	var docsInserted int64
	docArray := []interface{}{}
	switch docType := doc.(type) {
	case map[string]interface{}:
		docsInserted = 1
		docArray = []interface{}{docType}
	case []interface{}:
		docsInserted = int64(len(docType))
		docArray = docType
	default:
		return 0, fmt.Errorf("unknown document type %T", docType)
	}

	ch, ok := m.batchMapTableToChan[project][dbAlias][col] // get channel for specified table
	if !ok {
		logrus.Errorf("error converting insert request to batch request unable to find channel for database %s & collection %s", dbAlias, col)
		return 0, fmt.Errorf("cannot find channel for database %s & collection %s", dbAlias, col)
	}
	ch.request <- batchRequest{document: docArray, response: response}
	result := <-response
	return docsInserted, result.err
}
