package eventing

import (
	"fmt"
	"time"
)

func (m *Module) routineProcessIntents() {
	defer m.wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	// for t := range ticker.C {
	// 	m.processIntents(&t)
	// }
	for {
		select {
		case t := <-ticker.C:
			m.processIntents(&t)
		case <-m.stopChan:
			fmt.Println("goroutine routineProcessIntents stopped")
			return
		}
	}
}

func (m *Module) routineProcessStaged() {
	defer m.wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	// for t := range ticker.C {
	// 	m.processStagedEvents(&t)
	// }
	for {
		select {
		case t := <-ticker.C:
			m.processStagedEvents(&t)
		case <-m.stopChan:
			fmt.Println("goroutine routineProcessStaged stopped")
			return
		}
	}
}
