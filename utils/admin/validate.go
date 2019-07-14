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
func (v *validator) startValidation(id, account, key string, mode int) error {
	// Set validation status to active
	v.setActive(true)

	if err := v.registerSpaceCloud(id, account, key, mode); err != nil {
		return err
	}

	go func() {
		timer := time.Now()
		for {
			if err := v.routineRead(); err != nil {
				log.Println("Validate: Error -", err)
			}

			if !v.isActive() {
				return
			}

			// Sleep for 5 minutes before connecting again
			time.Sleep(5 * time.Minute)

			// Check if 15 days are lapsed without authorization
			if time.Since(timer).Hours() > 24*15 {

				// Stop the validation process
				v.stopValidation()
				v.reduceMode()
				return
			}

			if err := v.registerSpaceCloud(id, account, key, mode); err != nil {
				log.Println("Validate: Error -", err)
			} else {
				timer = time.Now()
			}
		}
	}()

	return nil
}

func (v *validator) stopValidation() {
	v.setActive(false)
	if v.socket != nil {
		v.socket.Close()
	}
}

func (v *validator) isActive() bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.active
}
func (v *validator) setActive(active bool) {
	v.lock.Lock()
	v.active = active
	v.lock.Unlock()
}
