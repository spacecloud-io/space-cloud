package validate

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (v *Validator) connect(url string) error {
	// Create a new client
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	v.socket = c

	// Start a routine to routinely send pings
	go func(socket *websocket.Conn) {
		ticker := time.NewTicker(20 * time.Second)
		for range ticker.C {
			if err := socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}(c)

	return nil
}

func (v *Validator) write(id, account, secret string) error {
	req := &model.RegisterRequest{ID: id, Account: account, Secret: secret}
	msg := &model.Message{Type: utils.TypeRegisterRequest, Data: req}
	return v.socket.WriteJSON(msg)
}

func (v *Validator) read() (msg *model.Message, err error) {
	msg = &model.Message{}
	err = v.socket.ReadJSON(msg)
	return
}

func (v *Validator) routineRead() error {
	for {
		msg, err := v.read()
		if err != nil {
			return err
		}

		switch msg.Type {
		case utils.TypeRegisterRequest:

			// For register request
			data := new(model.RegisterResponse)
			mapstructure.Decode(msg.Data, data)

			if !data.Ack {
				log.Println("Validate Error -", data.Error)
				os.Exit(-1)
			}

		case utils.TypeProjectFeed:
			// For project feed request
			data := new(model.ProjectFeed)
			mapstructure.Decode(msg.Data, data)

			// for project config
			conf := new(config.Project)
			mapstructure.Decode(data.Config, conf)

			err := v.projects.StoreProject(data.Project, conf)
			if err != nil {
				log.Println("Project Load Error -", err)
			}
		}
	}
}
