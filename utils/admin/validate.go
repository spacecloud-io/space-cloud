package admin

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const url = "wss://spaceuptech.com/v1/authenticate/socket/json"

// validator is the object which valiates the space cloud instance
type validator struct {
	lock       sync.Mutex
	socket     *websocket.Conn
	active     bool
	reduceMode func()
}

// New creates a new instance of validator
func newValidator(cb func()) *validator {
	return &validator{active: false, reduceMode: cb}
}

// Start starts the validation process
func (v *validator) startValidation(id, account, key string) error {
	// Set validation status to active
	v.lock.Lock()
	v.active = true
	v.lock.Unlock()

	if err := v.registerSpaceCloud(id, account, key); err != nil {
		return err
	}

	go func() {
		timer := time.Now()
		for {
			if !v.isActive() {
				return
			}

			if err := v.routineRead(); err != nil {
				log.Println("Validate: Error -", err)
			}

			// Sleep for 5 minutes before connecting again
			time.Sleep(5 * time.Minute)

			// Check if 15 days are lapsed without authorization
			if time.Since(timer).Hours() > 24*15 {

				// Reduce op mode to open source
				v.reduceMode()
				return
			}

			if err := v.registerSpaceCloud(id, account, key); err != nil {
				log.Println("Validate: Error -", err)
			} else {
				timer = time.Now()
			}
		}
	}()

	return nil
}

func (v *validator) stopValidation() {
	v.lock.Lock()
	v.active = false
	v.lock.Unlock()
}

func (v *validator) isActive() bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.active
}
