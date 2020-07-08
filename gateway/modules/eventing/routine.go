package eventing

func (m *Module) routineProcessIntents() {
	for t := range m.ticker.C {
		m.processIntents(&t)
	}
}

func (m *Module) routineProcessStaged() {
	for t := range m.ticker.C {
		m.processStagedEvents(&t)
	}
}
