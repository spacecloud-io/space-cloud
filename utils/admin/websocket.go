package admin

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (v *validator) connect(url string) error {
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
				log.Println("Validator Error:", err)
				return
			}
		}
	}(c)

	return nil
}

func (v *validator) write(id, userID, key string, mode int) error {
	req := &model.RegisterRequest{ID: id, UserID: userID, Key: key, Mode: mode}
	msg := &model.Message{Type: utils.TypeRegisterRequest, Data: req}
	return v.socket.WriteJSON(msg)
}

func (v *validator) read() (msg *model.Message, err error) {
	msg = &model.Message{}
	err = v.socket.ReadJSON(msg)
	return
}

func (v *validator) routineRead() error {
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
				// Reduce op mode to open source
				v.stopValidation()
			}
		}
	}
}

func (v *validator) registerSpaceCloud(id, userID, secret string, mode int) error {
	err := v.connect(url)
	if err != nil {
		return err
	}

	err = v.write(id, userID, secret, mode)
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
