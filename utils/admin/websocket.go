package admin

import (
	"errors"
	"log"
	"os"
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

func (v *validator) write(id, account, key string) error {
	req := &model.RegisterRequest{ID: id, Account: account, Key: key}
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
				os.Exit(-1)
			}
		}
	}
}

func (v *validator) registerSpaceCloud(id, account, secret string) error {
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
