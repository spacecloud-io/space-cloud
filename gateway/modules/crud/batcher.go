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
	batchMap map[string]map[string]batchChannels // key here is dbAlias & it's value's key is table name

	batchChannels struct {
		request  batchRequestChan
		response chan error
		close    chan bool
	}

	batchRequestChan chan batchRequest

	batchRequest struct {
		document interface{}
		project  string
	}
)

// CloseBatchOperation closes all go routines associated with individual collection for batch operation
func (m *Module) CloseBatchOperation() {
	m.Lock()
	defer m.Unlock()
	for _, value := range m.batchMapTableToChan {
		for _, channels := range value {
			close(channels.close)
		}
	}
}

// initBatchOperation creates go routines for executing batch operation associated with individual collection
func (m *Module) initBatchOperation(crud config.Crud) {
	batch := batchMap{}
	// log.Println("go routines before creating executors -", runtime.NumGoroutine())
	for dbAlias, dbInfo := range crud {
		if dbInfo.Enabled {
			for tableName := range dbInfo.Collections {
				done := make(chan bool)                      // channel for closing go routine
				addInsertToBatchCh := make(batchRequestChan) // channel for adding request to batch op
				response := make(chan error)                 // channel for sending op response back to client
				go m.insertBatchExecutor(done, response, addInsertToBatchCh, dbInfo.BatchTime, dbAlias, tableName)
				if batch[dbAlias] == nil {
					batch[dbAlias] = map[string]batchChannels{tableName: {request: addInsertToBatchCh, response: response, close: done}}
					continue
				}
				batch[dbAlias][tableName] = batchChannels{request: addInsertToBatchCh, response: response, close: done}
			}
		}
	}
	m.batchMapTableToChan = batch
	// log.Println("go routines after creating executors -", runtime.NumGoroutine())
}

func (m *Module) insertBatchExecutor(done chan bool, response chan error, addInsertToBatchCh batchRequestChan, batchTime int, dbAlias, tableName string) {
	project := ""
	batchRequests := make([]interface{}, 0)
	if batchTime <= 0 { // when new project is created set default time to 200 milli seconds
		batchTime = 200
	}
	ticker := time.NewTicker(time.Duration(batchTime) * time.Millisecond)
	for {
		select {
		case <-done:
			ticker.Stop()
			logrus.Debugf("closing batcher for database %s table %s", dbAlias, tableName)
			return
		case v := <-addInsertToBatchCh:
			project = v.project
			batchRequests = append(batchRequests, v.document.([]interface{}))
			// log.Println("added batch request", batchRequests)
		case <-ticker.C:
			if len(batchRequests) != 0 {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO CHECK CONTEXT TIME
				if err := m.Create(ctx, dbAlias, project, tableName, &model.CreateRequest{Operation: utils.All, Document: batchRequests}); err != nil {
					logrus.Errorf("error executing batch request for database %s table %s - %s", dbAlias, tableName, err)
					response <- err
				}
				cancel()                               // close context
				batchRequests = make([]interface{}, 0) // clear the requests array
				response <- nil
			}
		}
	}
}
