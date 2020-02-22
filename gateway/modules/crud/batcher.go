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
		close   chan struct{}
	}

	batchRequestChan chan batchRequest

	batchRequest struct {
		document interface{}
		response chan error
	}
)

// CloseBatchOperation closes all go routines associated with individual collection for batch operation
func (m *Module) CloseBatchOperation() {
	for _, projectValue := range m.batchMapTableToChan {
		for _, dbAliasValue := range projectValue {
			for _, tableValue := range dbAliasValue {
				tableValue.close <- struct{}{}
			}
		}
	}
}

// initBatchOperation creates go routines for executing batch operation associated with individual collection
func (m *Module) initBatchOperation(project string, crud config.Crud) {
	batch := batchMap{}
	for dbAlias, dbInfo := range crud {
		if dbInfo.Enabled {
			for tableName, tableInfo := range dbInfo.Collections {
				done := make(chan struct{})                                          // channel for closing go routine
				addInsertToBatchCh := make(batchRequestChan, tableInfo.BatchRecords) // channel for adding request to batch op
				go m.insertBatchExecutor(done, addInsertToBatchCh, dbInfo.BatchTime, project, dbAlias, tableName, tableInfo.BatchRecords)
				if batch[project] == nil {
					batch[project] = map[string]map[string]batchChannels{dbAlias: {tableName: {request: addInsertToBatchCh, close: done}}}
					continue
				}
				batch[project][dbAlias][tableName] = batchChannels{request: addInsertToBatchCh, close: done}
			}
		}
	}
	m.batchMapTableToChan = batch
}

func (m *Module) insertBatchExecutor(done chan struct{}, addInsertToBatchCh batchRequestChan, batchTime int, project, dbAlias, tableName string, batchRecordLimit int) {
	noOfRequests := 0
	responseChannels := make([]chan error, 0)
	batchRequests := make([]interface{}, 0)
	if batchTime <= 0 { // when new project is created set default time to 200 milli seconds
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
			close(addInsertToBatchCh)
			close(done)
			logrus.Debugf("closing batcher for database %s table %s", dbAlias, tableName)
			return
		case v := <-addInsertToBatchCh:
			noOfRequests++
			responseChannels = append(responseChannels, v.response)
			switch docType := v.document.(type) {
			case map[string]interface{}:
				batchRequests = append(batchRequests, docType)
			case []interface{}:
				batchRequests = append(batchRequests, docType...)
			}
			if noOfRequests == batchRecordLimit {
				m.executeBatch(project, dbAlias, tableName, batchRequests, responseChannels)
				batchRequests = make([]interface{}, 0) // clear the requests array
				noOfRequests = 0
				responseChannels = make([]chan error, 0) // clear the response channels array
				// reset ticker
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(batchTime) * time.Millisecond)
			}
		case <-ticker.C:
			if len(batchRequests) != 0 {
				m.executeBatch(project, dbAlias, tableName, batchRequests, responseChannels)
				batchRequests = make([]interface{}, 0) // clear the requests array
				noOfRequests = 0
				responseChannels = make([]chan error, 0) // clear the response channels array
			}
		}
	}
}

func (m *Module) executeBatch(project, dbAlias, tableName string, batchRequests []interface{}, responseChannels []chan error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := m.Create(ctx, dbAlias, project, tableName, &model.CreateRequest{Operation: utils.All, Document: batchRequests}); err != nil {
		logrus.Errorf("error executing batch request for database %s table %s - %s", dbAlias, tableName, err)
		for _, responseChan := range responseChannels {
			responseChan <- err
		}
		return
	}
	// send response to all client request
	for _, responseChan := range responseChannels {
		responseChan <- nil
	}
}
