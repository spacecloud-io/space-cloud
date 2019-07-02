package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
)

// GRPCServiceClient is the object handling all client interactions
type GRPCServiceClient struct {
	id           string
	channel      chan *model.Message
	ctx          context.Context
	cancel       context.CancelFunc
	streamServer proto.SpaceCloud_ServiceServer
}

// RoutineWrite starts a json writer routine
func (c *GRPCServiceClient) RoutineWrite() {
	for res := range c.channel {
		switch res.Type {
		case utils.TypeServiceRequest:
			reqMsg, ok := res.Data.(*model.FunctionsPayload)
			if !ok {
				break
			}

			authData, _ := json.Marshal(reqMsg.Auth)
			paramsData, _ := json.Marshal(reqMsg.Params)
			c.streamServer.Send(&proto.FunctionsPayload{
				Id:       reqMsg.ID,
				Type:     utils.TypeServiceRequest,
				Auth:     authData,
				Params:   paramsData,
				Function: reqMsg.Func,
			})

		case utils.TypeServiceRegister:
			reqMsg, ok := res.Data.(map[string]interface{})
			if !ok {
				log.Println("GRPC Service Error - Invalid data type", res.Data)
				break
			}
			paramsData, _ := json.Marshal(reqMsg)
			c.streamServer.Send(&proto.FunctionsPayload{
				Id:     res.ID,
				Type:   utils.TypeServiceRegister,
				Params: paramsData,
			})

		default:
			log.Println("GRPC Service Error - Invalid request type", res.Type)
		}
	}

}

// Write wrties the object to the client
func (c *GRPCServiceClient) Write(res *model.Message) {
	select {
	case <-c.ctx.Done():
	case c.channel <- res:
	}
}

// Close closes the client
func (c *GRPCServiceClient) Close() {
	c.cancel()
	close(c.channel)
}

// Read startes a blocking reader routine
func (c *GRPCServiceClient) Read(cb DataCallback) {
	for {
		in, err := c.streamServer.Recv()
		if err != nil {
			if err != nil {
				log.Println("GRPC Service Receive Error -", err)
				return
			}
		}

		switch in.Type {
		case utils.TypeServiceRegister:
			data := map[string]interface{}{"service": in.Service, "token": in.Token, "project": in.Project}
			msg := &model.Message{ID: in.Id, Type: utils.TypeServiceRegister, Data: data}

			// Close the reader if callback returned false
			next := cb(msg)
			if !next {
				return
			}

		case utils.TypeServiceRequest:
			var params interface{}
			json.Unmarshal(in.Params, &params)
			data := map[string]interface{}{
				"params": params,
				"id":     in.Id,
				"error":  in.Error,
			}
			msg := &model.Message{ID: in.Id, Type: utils.TypeServiceRequest, Data: data}

			// Close the reader if callback returned false
			next := cb(msg)
			if !next {
				return
			}

		default:
			log.Println("GRPC Service Error - Invalid request type", in.Type)
		}
	}

}

// ClientID returns the client's id
func (c *GRPCServiceClient) ClientID() string {
	return c.id
}

// Context returns the client's context
func (c *GRPCServiceClient) Context() context.Context {
	return c.ctx
}
