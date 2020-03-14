package crud

import (
	"context"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/gateway/config"
)

// data structure for storing a group of channels
// against particular table/collection of particular database
type (
	batchMap map[string]map[string]map[string]batchChannels // keys are project name, dbAlias, table name

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

// closeBatchOperation closes all go routines associated with individual collection for batch operation
func (m *Module) closeBatchOperation() {
	for _, projectValue := range m.batchMapTableToChan {
		for _, dbAliasValue := range projectValue {
			for _, tableValue := range dbAliasValue {
				tableValue.closeC <- struct{}{}
			}
		}
	}
}

// initBatchOperation creates go routines for executing batch operation associated with individual collection
func (m *Module) initBatchOperation(project string, crud config.Crud) {
	batch := batchMap{}
	for dbAlias, dbInfo := range crud {
		if dbInfo.Enabled {
			for tableName := range dbInfo.Collections {
				closeC := make(chan struct{})                    // channel for closing go routine
				addInsertToBatchCh := make(batchRequestChan, 20) // channel for adding request to batch op
				go m.insertBatchExecutor(closeC, addInsertToBatchCh, dbInfo.BatchTime, project, dbAlias, tableName, dbInfo.BatchRecords)
				if batch[project] == nil {
					batch[project] = map[string]map[string]batchChannels{dbAlias: {tableName: {request: addInsertToBatchCh, closeC: closeC}}}
					continue
				}
				if batch[project][dbAlias] == nil {
					batch[project][dbAlias] = map[string]batchChannels{tableName: {request: addInsertToBatchCh, closeC: closeC}}
					continue
				}
				batch[project][dbAlias][tableName] = batchChannels{request: addInsertToBatchCh, closeC: closeC}
			}
		}
	}
	m.batchMapTableToChan = batch
}

func (m *Module) insertBatchExecutor(done chan struct{}, addInsertToBatchCh batchRequestChan, batchTime int, project, dbAlias, tableName string, batchRecordLimit int) {
	responseChannels := make([]batchResponseChan, 0)
	batchRequests := make([]interface{}, 0)
	if batchTime <= 0 { // when new project is created set default time to 100 milli seconds
		batchTime = 100
	}
	if batchRecordLimit <= 0 {
		batchRecordLimit = 100 // when new project is created set default batch record limit to 100
	}
	ticker := time.NewTicker(time.Duration(batchTime) * time.Millisecond)
	for {
		select {
		case <-done:
			ticker.Stop()
			// safe operation since SetConfig will hold a lock preventing others from writing into this channel after its closed
			close(addInsertToBatchCh)
			close(done)
			logrus.Debugf("closing batcher for database %s table %s", dbAlias, tableName)
			return
		case v := <-addInsertToBatchCh:
			responseChannels = append(responseChannels, v.response)
			batchRequests = append(batchRequests, v.documents...)
			if len(batchRequests) >= batchRecordLimit {
				m.executeBatch(project, dbAlias, tableName, batchRequests, responseChannels)
				batchRequests = make([]interface{}, 0)          // clear the requests array
				responseChannels = make([]batchResponseChan, 0) // clear the response channels array
				// reset ticker
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(batchTime) * time.Millisecond)
			}
		case <-ticker.C:
			if len(batchRequests) > 0 {
				m.executeBatch(project, dbAlias, tableName, batchRequests, responseChannels)
				batchRequests = make([]interface{}, 0)          // clear the requests array
				responseChannels = make([]batchResponseChan, 0) // clear the response channels array
			}
		}
	}
}

func (m *Module) executeBatch(project, dbAlias, tableName string, batchRequests []interface{}, responseChannels []batchResponseChan) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := m.InternalCreate(ctx, dbAlias, project, tableName, &model.CreateRequest{Operation: utils.All, Document: batchRequests}, true); err != nil {
		logrus.Errorf("error executing batch request for database %s table %s - %s", dbAlias, tableName, err)
		m.sendResponses(responseChannels, batchResponse{err: err})
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
