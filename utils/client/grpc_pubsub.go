package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
)

// GRPCPubsubClient is the object handling all client interactions
type GRPCPubsubClient struct {
	id           string
	channel      chan *model.Message
	ctx          context.Context
	cancel       context.CancelFunc
	streamServer proto.SpaceCloud_PubsubSubscribeServer
}

// RoutineWrite starts a json writer routine
func (c *GRPCPubsubClient) RoutineWrite() {
	for res := range c.channel {
		switch res.Type {
		case utils.TypePubsubSubscribeFeed:
			msg, ok := res.Data.(*model.PubsubMsg)
			if !ok {
				break
			}
			b, err := json.Marshal(msg)
			if err != nil {
				log.Println("GRPC Pubsub Error - Could not marshal", res.Data)
				break
			}

			c.streamServer.Send(&proto.PubsubMsgResponse{
				Id:   res.ID,
				Type: utils.TypePubsubSubscribeFeed,
				Msg:  b,
			})

		case utils.TypePubsubSubscribe, utils.TypePubsubUnsubscribe, utils.TypePubsubUnsubscribeAll:
			resp, ok := res.Data.(model.PubsubMsgResponse)
			if !ok {
				log.Println("GRPC Pubsub Error - Invalid data type", res.Data)
				break
			}
			c.streamServer.Send(&proto.PubsubMsgResponse{
				Id:     res.ID,
				Type:   res.Type,
				Status: resp.Status,
				Error:  resp.Error,
			})

		default:
			log.Println("GRPC Pubsub Error - Invalid request type", res.Type)
		}
	}

}

// Write wrties the object to the client
func (c *GRPCPubsubClient) Write(res *model.Message) {
	select {
	case <-c.ctx.Done():
	case c.channel <- res:
	}
}

// Close closes the client
func (c *GRPCPubsubClient) Close() {
	c.cancel()
	close(c.channel)
}

// Read startes a blocking reader routine
func (c *GRPCPubsubClient) Read(cb DataCallback) {
	for {
		in, err := c.streamServer.Recv()
		if err != nil {
			if err != nil {
				log.Println("GRPC Pubsub Receive Error -", err)
				return
			}
		}

		switch in.Type {
		case utils.TypePubsubSubscribe, utils.TypePubsubUnsubscribe, utils.TypePubsubUnsubscribeAll:
			data := map[string]interface{}{"subject": in.Subject, "queue": in.Queue, "type": in.Type, "token":in.Token, "project":in.Project, "id":in.Id}
			msg := &model.Message{ID: in.Id, Type: in.Type, Data: data}

			// Close the reader if callback returned false
			next := cb(msg)
			if !next {
				return
			}

		default:
			log.Println("GRPC Pubsub Error - Invalid request type", in.Type)
		}
	}

}

// ClientID returns the client's id
func (c *GRPCPubsubClient) ClientID() string {
	return c.id
}

// Context returns the client's context
func (c *GRPCPubsubClient) Context() context.Context {
	return c.ctx
}
