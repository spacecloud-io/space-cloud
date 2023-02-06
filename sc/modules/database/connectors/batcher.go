package connectors

import (
	"context"
	"fmt"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

// data structure for storing a group of channels
// against particular table/collection of particular database
type (
	batchMap map[string]batchChannels // key is table name

	batchChannels struct {
		request batchRequestChan
		closeC  chan struct{}
	}

	batchRequestChan  chan batchRequest
	batchResponseChan chan batchResponse

	batchRequest struct {
		documents []interface{}
		response  batchResponseChan
	}
	batchResponse struct {
		err error
	}
)

func (m *Module) createBatch(ctx context.Context, col string, doc interface{}) (int64, error) {
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

	ch, ok := m.batchMapTableToChan[col] // get channel for specified table
	if !ok {
		return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot convert insert request to batch request", fmt.Errorf("cannot find channel for database %s & collection %s", m.dbConfig.DbAlias, col), nil)
	}
	ch.request <- batchRequest{documents: docArray, response: response}
	result := <-response
	return docsInserted, result.err
}

// closeBatchOperation closes all go routines associated with individual collection for batch operation
func (m *Module) closeBatchOperation() {
	for _, tableValue := range m.batchMapTableToChan {
		tableValue.closeC <- struct{}{}
	}
}

// initBatchOperation creates go routines for executing batch operation associated with individual collection
func (m *Module) initBatchOperation() {
	batch := batchMap{}
	for _, schema := range m.dbSchemas {

		closeC := make(chan struct{})                    // channel for closing go routine
		addInsertToBatchCh := make(batchRequestChan, 20) // channel for adding request to batch op
		go m.insertBatchExecutor(closeC, addInsertToBatchCh, m.dbConfig.BatchTime, m.dbConfig.DbAlias, schema.Table, m.dbConfig.BatchRecords)
		batch[schema.Table] = batchChannels{request: addInsertToBatchCh, closeC: closeC}
	}
	m.batchMapTableToChan = batch
}

func (m *Module) insertBatchExecutor(done chan struct{}, addInsertToBatchCh batchRequestChan, batchTime int, dbAlias, tableName string, batchRecordLimit int) {
	responseChannels := make([]batchResponseChan, 0)
	batchRequests := make([]interface{}, 0)
	if batchTime <= 0 { // when new project is created set default time to 100 milli seconds
		batchTime = 200
	}
	if batchRecordLimit <= 0 {
		batchRecordLimit = 200 // when new project is created set default batch record limit to 100
	}
	ticker := time.NewTicker(time.Duration(batchTime) * time.Millisecond)
	for {
		select {
		case <-done:
			ticker.Stop()
			// safe operation since SetConfig will hold a lock preventing others from writing into this channel after its closed
			close(addInsertToBatchCh)
			close(done)
			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("closing batcher for database %s table %s", dbAlias, tableName), nil)
			return
		case v := <-addInsertToBatchCh:
			responseChannels = append(responseChannels, v.response)
			batchRequests = append(batchRequests, v.documents...)
			if len(batchRequests) >= batchRecordLimit {
				m.executeBatch(dbAlias, tableName, batchRequests, responseChannels)
				batchRequests = make([]interface{}, 0)          // clear the requests array
				responseChannels = make([]batchResponseChan, 0) // clear the response channels array
				// reset ticker
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(batchTime) * time.Millisecond)
			}
		case <-ticker.C:
			if len(batchRequests) > 0 {
				m.executeBatch(dbAlias, tableName, batchRequests, responseChannels)
				batchRequests = make([]interface{}, 0)          // clear the requests array
				responseChannels = make([]batchResponseChan, 0) // clear the response channels array
			}
		}
	}
}

func (m *Module) executeBatch(dbAlias, tableName string, batchRequests []interface{}, responseChannels []batchResponseChan) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.connector.IsClientSafe(ctx); err != nil {
		m.sendResponses(responseChannels, batchResponse{err: helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error executing batch request for database %s table %s", dbAlias, tableName), err, nil)})
	}

	if _, err := m.connector.Create(ctx, tableName, &model.CreateRequest{Operation: utils.All, Document: batchRequests}); err != nil {
		m.sendResponses(responseChannels, batchResponse{err: helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error executing batch request for database %s table %s", dbAlias, tableName), err, nil)})
		return
	}
	// send response to all client request
	m.sendResponses(responseChannels, batchResponse{err: nil})
}

func (m *Module) sendResponses(responseChannels []batchResponseChan, response batchResponse) {
	for _, responseChan := range responseChannels {
		responseChan <- response
	}
}
