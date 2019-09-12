package eventing

import (
	"log"

	nats "github.com/nats-io/nats.go"
)

func (m *Module) routineEvents(channel chan *nats.Msg) {
	for msg := range channel {
		log.Println("Eventing received:", string(msg.Data))
	}
}
