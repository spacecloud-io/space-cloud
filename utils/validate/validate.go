package validate

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

const url = "wss://spaceuptech.com/v1/authenticate/socket/json"

// Validator is the object which valiates the space cloud instance
type Validator struct {
	socket   *websocket.Conn
	projects *projects.Projects
}

// New creates a new instance of validator
func New(p *projects.Projects) *Validator {
	return &Validator{projects: p}
}

// Start starts the validation process
func (v *Validator) Start(id, account, secret string) error {
	if err := v.registerSpaceCloud(id, account, secret); err != nil {
		return err
	}

	go func() {
		for {
			err := v.routineRead()
			if err != nil {
				log.Println("Validate: Error -", err)
			}

			// Sleep for 5 minutes before connecting again
			time.Sleep(5 * time.Minute)

			err = v.registerSpaceCloud(id, account, secret)
			if err != nil {
				log.Println("Validate: Error -", err)
			}
		}
	}()

	return nil
}

func (v *Validator) registerSpaceCloud(id, account, secret string) error {
	err := v.connect(url)
	if err != nil {
		return err
	}

	err = v.write(id, account, secret)
	if err != nil {
		return err
	}

	msg, err := v.read()
	if err != nil {
		return err
	}

	if msg.Type == utils.TypeRegisterRequest {
		// For register request
		data := new(model.RegisterResponse)
		mapstructure.Decode(msg.Data, data)

		if !data.Ack {
			return errors.New(data.Error)
		}
	}

	return nil
}
