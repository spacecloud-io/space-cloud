package eventing

import (
	"time"
)

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
