package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/pb"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
)

type counter struct {
	value, nos int32
}

type metric struct {
	Value int32 `json:"val"`
	Ts    int64 `json:"ts"`
}

func (s *Server) aggregate() {
	// Create a 60s aggregator and a 6s aggregator
	a60 := newAggregator()
	a6 := newAggregator()

	// Take the current time snapshot
	now := time.Now().Unix()

	// Create stream
	stream := s.db.NewStream()
	stream.NumGo = 16
	stream.Prefix = []byte("metrics")
	stream.Send = func(list *pb.KVList) error {
		for _, kv := range list.Kv {
			// Get the project id, service, version and node id
			array := strings.Split(string(kv.Key), "/")
			project, service, version, nodeID := array[1], array[2], array[3], array[4]

			// Unmarshal the metrics from badger
			m := new(metric)
			_ = json.Unmarshal(kv.Value, m)

			// Add the metric to the 60s aggregator. Add it to the 6s aggregator only if its less that 6s old.
			a60.add(project, service, version, nodeID, m.Value)
			if m.Ts+6 >= now {
				a6.add(project, service, version, nodeID, m.Value)
			}
		}
		return nil
	}

	// Orchestrate the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		logrus.Errorln("Could start stream from badger:", err)
		return
	}

	// Services that require scale adjusting

	// Iterate over all 60 second aggregations
	a60.iterate(func(project, service, version string, value int32) {

		// Enter panic mode if 6 second average is twice or half the value of 60 second average. In panic mode, we make all decision based on the
		// count of the 6 second average
		v6 := a6.get(project, service, version)
		if v6 != 0 && (v6 >= value*2 || v6 <= value/2) {
			value = v6
		}

		// Adjust the scale of the service
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			if err := s.driver.AdjustScale(ctx, &model.Service{ProjectID: project, ID: service, Version: version}, value); err != nil {
				logrus.Errorf("Could not adjust scale of service (%s:%s): %s", project, service, err.Error())
			}
		}()

		a6.delete(project, service, version)
	})

	a6.iterate(func(project, service, version string, value int32) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			if err := s.driver.AdjustScale(ctx, &model.Service{ProjectID: project, ID: service, Version: version}, value); err != nil {
				logrus.Errorf("Could not adjust scale of service (%s:%s): %s", project, service, err.Error())
			}
		}()
	})
}

func (s *Server) routineAdjustScale() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		s.aggregate()
	}
}

func (s *Server) routineDumpDetails() {
	messages := make([]*model.ProxyMessage, 0)
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if len(messages) > 0 {
				if err := s.flushMetrics(messages); err != nil {
					logrus.Errorln("Could not flush metrics to disk:", err)
				}
				messages = []*model.ProxyMessage{}
			}
		case msg := <-s.chAppend:
			messages = append(messages, msg)
		}
	}
}

func (s *Server) flushMetrics(metrics []*model.ProxyMessage) error {
	return s.db.Update(func(txn *badger.Txn) error {
		for _, m := range metrics {
			// Prepare the key and values
			key := fmt.Sprintf("metrics/%s/%s/%s/%s/%s", m.Project, m.Service, m.Version, m.NodeID, ksuid.New().String())
			data, _ := json.Marshal(&metric{Ts: time.Now().Unix(), Value: m.ActiveRequests})
			// Set entry in badger
			e := badger.NewEntry([]byte(key), data).WithTTL(time.Minute)
			if err := txn.SetEntry(e); err != nil {
				return err
			}
		}

		return nil
	})
}
