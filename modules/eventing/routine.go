package eventing

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/model"
)

func (m *Module) initEventWorkers(channel chan *nats.Msg, workers int) {
	for i := 0; i < workers; i++ {
		go m.routineEvents(channel)
	}
}

func (m *Module) routineEvents(channel chan *nats.Msg) {
	for msg := range channel {
		var eventDocs []*model.EventDocument
		if err := json.Unmarshal(msg.Data, &eventDocs); err != nil {
			log.Println("Eventing: Unable to unmarshal event documents -", err)
			continue
		}

		m.processBroadcastEvents(eventDocs)
	}
}

func (m *Module) routineProcessIntents() {
	ticker := time.NewTicker(10 * time.Second)
	for t := range ticker.C {
		m.processIntents(&t)
	}
}

func (m *Module) routineProcessStaged() {
	ticker := time.NewTicker(10 * time.Second)
	for t := range ticker.C {
		m.processStagedEvents(&t)
	}
}
