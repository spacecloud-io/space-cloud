package crud

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"
)

func (m *Module) createBatch(ctx context.Context, project, dbAlias, col string, doc interface{}) (int64, error) {
	response := make(batchResponseChan, 1)
	defer close(response)

	var docsInserted int64
	var docArray []interface{}
	switch docType := doc.(type) {
	case map[string]interface{}:
		docsInserted = 1
		docArray = []interface{}{docType}
	case []interface{}:
		docsInserted = int64(len(docType))
		docArray = docType
	default:
		return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Cannot create batch request unkownd doc type (%T) provided)", docType), nil, nil)
	}

	// Simply return if 0 docs are to be inserted
	if docsInserted == 0 {
		return 0, nil
	}

	ch, ok := m.batchMapTableToChan[project][dbAlias][col] // get channel for specified table
	if !ok {
		return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot convert insert request to batch request", fmt.Errorf("cannot find channel for database %s & collection %s", dbAlias, col), nil)
	}
	ch.request <- batchRequest{documents: docArray, response: response}
	result := <-response
	return docsInserted, result.err
}

func getPreparedQueryKey(dbAlias, id string) string {
	return fmt.Sprintf("%s--%s", dbAlias, id)
}
