package eventing

func (m *Module) routineProcessIntents() {
	for {
		select {
		case <-m.done:
			return
		case t := <-m.ticker.C:
			m.processIntents(&t)
		}
	}
}

func (m *Module) routineProcessStaged() {
	for {
		select {
		case <-m.done:
			return
		case t := <-m.ticker.C:
			m.processStagedEvents(&t)
		}
	}
}
