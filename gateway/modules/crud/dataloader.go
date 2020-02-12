package crud

import (
	"context"
	"sync"

	"github.com/graph-gophers/dataloader"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type resultsHolder struct {
	sync.Mutex
	results      []*dataloader.Result
	whereClauses []interface{}
}

func (holder *resultsHolder) getResults() []*dataloader.Result {
	holder.Lock()
	defer holder.Unlock()

	return holder.results
}

func (holder *resultsHolder) addResult(i int, result *dataloader.Result) {
	holder.Lock()
	holder.results[i] = result
	holder.Unlock()
}

func (holder *resultsHolder) getWhereClauses() []interface{} {
	holder.Lock()
	defer holder.Unlock()

	return holder.whereClauses
}

func (holder *resultsHolder) addWhereClause(whereClause map[string]interface{}) {
	holder.Lock()
	holder.whereClauses = append(holder.whereClauses, whereClause)
	holder.Unlock()
}

func (holder *resultsHolder) fillResults(res []interface{}) {
	holder.Lock()
	defer holder.Unlock()

	// Create a where clause index
	index := 0

	length := len(holder.results)
	for i := 0; i < length; i++ {

		// Continue if result already has a value
		if holder.results[i] != nil {
			continue
		}

		// Get the where clause
		whereClause := holder.whereClauses[index]

		docs := []interface{}{}
		for _, doc := range res {
			if utils.Validate(whereClause.(map[string]interface{}), doc) {
				docs = append(docs, doc)
			}
		}

		// Increment the where clause index
		index++

		// Store the matched docs in result
		holder.results[i] = &dataloader.Result{Data: docs}
	}
}

func (holder *resultsHolder) fillErrorMessage(err error) {
	holder.Lock()

	length := len(holder.results)
	for i := 0; i < length; i++ {
		if holder.results[i] == nil {
			holder.results[i] = &dataloader.Result{Error: err}
		}
	}
	holder.Unlock()
}

type loaderMap struct {
	lock sync.Mutex
	m    map[string]*dataloader.Loader
}

func newLoaderMap() *loaderMap {
	return &loaderMap{m: map[string]*dataloader.Loader{}}
}

func (l *loaderMap) get(key string, m *Module) *dataloader.Loader {
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, ok := l.m[key]; !ok {
		l.m[key] = m.createLoader()
	}

	return l.m[key]
}

func (m *Module) createLoader() *dataloader.Loader {
	// DataLoaderBatchFn is the batch function of the data loader
	cache := &dataloader.NoCache{}
	return dataloader.NewBatchedLoader(m.dataLoaderBatchFn, dataloader.WithCache(cache))
}

func (m *Module) dataLoaderBatchFn(c context.Context, keys dataloader.Keys) []*dataloader.Result {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	var dbType, col string

	// Return if there are no keys
	if len(keys) == 0 {
		return []*dataloader.Result{}
	}

	holder := resultsHolder{
		results:      make([]*dataloader.Result, len(keys)),
		whereClauses: []interface{}{},
	}

	for index, key := range keys {
		req := key.(model.ReadRequestKey)

		dbType = req.DBType
		col = req.Col

		// Execute query immediately if it has options
		if req.HasOptions {
			// Add task to wait group
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				crud, err := m.getCrudBlock(dbType)
				if err != nil {
					cancel()
					holder.addResult(i, &dataloader.Result{Error: err})
					return
				}

				if err := crud.IsClientSafe(); err != nil {
					cancel()
					holder.addResult(i, &dataloader.Result{Error: err})
					return
				}
				// Execute the query
				_, res, err := crud.Read(ctx, m.project, req.Col, &req.Req)
				if err != nil {

					// Cancel the context and add the error response to the result
					cancel()
					holder.addResult(i, &dataloader.Result{Error: err})
					return
				}

				// Add the response to the result
				holder.addResult(i, &dataloader.Result{Data: res})
			}(index)

			// Continue to the next key
			continue
		}

		// Append the where clause to the list
		holder.addWhereClause(req.Req.Find)
	}

	// Wait for all results to be done
	wg.Wait()

	// Prepare a merged request
	req := model.ReadRequest{Find: map[string]interface{}{"$or": holder.getWhereClauses()}, Operation: utils.All, Options: &model.ReadOptions{}}

	// Fire the merged request

	crud, err := m.getCrudBlock(dbType)
	if err != nil {
		return []*dataloader.Result{}
	}

	if err := crud.IsClientSafe(); err != nil {
		return []*dataloader.Result{}
	}
	_, res, err := crud.Read(ctx, m.project, col, &req)
	if err != nil {
		holder.fillErrorMessage(err)
	} else {
		holder.fillResults(res.([]interface{}))
	}

	// do some async work to get data for specified keys
	// append to this list resolved values
	return holder.getResults()
}
