package eventing

import "time"

func (m *Module) routineProcessIntents() {
	m.tickerIntent = time.NewTicker(10 * time.Second)
	for t := range m.tickerIntent.C {
		m.processIntents(&t)
	}
}

func (m *Module) routineProcessStaged() {
	m.tickerStaged = time.NewTicker(10 * time.Second)
	for t := range m.tickerStaged.C {
		m.processStagedEvents(&t)
	}
}
