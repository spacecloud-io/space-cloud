package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
)

// GRPCRealtimeClient is the object handling all client interactions
type GRPCRealtimeClient struct {
	id      string
	channel chan *model.Message
	ctx     context.Context
	cancel  context.CancelFunc
	stream  proto.SpaceCloud_RealTimeServer
}

// RoutineWrite starts a json writer routine
func (c *GRPCRealtimeClient) RoutineWrite() {
	for res := range c.channel {
		//Convert the Message into RealTime response.
		switch res.Type {
		case utils.TypeRealtimeSubscribe, utils.TypeRealtimeUnsubscribe:
			//Decode the Message
			responseMsg := res.Data.(model.RealtimeResponse)
			feedData := make([]*proto.FeedData, len(responseMsg.Docs))
			for i, feed := range responseMsg.Docs {
				payload, err := json.Marshal(feed.Payload)
				if err != nil {
					log.Println(err)
					return
				}
				feedData[i] = &proto.FeedData{QueryId: feed.QueryID, DocId: feed.DocID, Type: feed.Type, Group: feed.Group, DbType: feed.DBType, Payload: payload, TimeStamp: feed.TimeStamp}
			}
			grpcResponse := proto.RealTimeResponse{Id: res.ID, Group: responseMsg.Group, Ack: responseMsg.Ack, Error: responseMsg.Error, FeedData: feedData}
			c.stream.Send(&grpcResponse)

		case utils.TypeRealtimeFeed:
			feed := res.Data.(*model.FeedData)
			feedData := make([]*proto.FeedData, 1)
			payload, err := json.Marshal(feed.Payload)
			if err != nil {
				log.Println(err)
				return
			}
			feedData[0] = &proto.FeedData{
				QueryId: feed.QueryID, DocId: feed.DocID, Type: feed.Type, Group: feed.Group, DbType: feed.DBType, Payload: payload, TimeStamp: feed.TimeStamp}
			grpcResponse := proto.RealTimeResponse{Id: res.ID, Group: res.Data.(*model.FeedData).Group, FeedData: feedData, Ack: true}
			c.stream.Send(&grpcResponse)
		}
	}
}

// Write wrties the object to the client
func (c *GRPCRealtimeClient) Write(res *model.Message) {
	select {
	case c.channel <- res:
	case <-c.ctx.Done():
	}
}

// Close closes the client
func (c *GRPCRealtimeClient) Close() {
	c.cancel()
	close(c.channel)
}

// Read startes a blocking reader routine
func (c *GRPCRealtimeClient) Read(cb DataCallback) {
	for {
		in, err := c.stream.Recv()
		if err != nil {
			log.Println("GRPC Error -", err)
			return
		}
		data := make(map[string]interface{})
		data["Token"] = in.Token
		data["DBType"] = in.DbType
		data["Project"] = in.Project
		data["Group"] = in.Group
		data["Type"] = in.Type
		data["ID"] = in.Id
		var temp interface{}
		err = json.Unmarshal(in.Where, &temp)
		if err != nil {
			log.Println(err)
			return
		}
		data["Where"] = temp

		msg := &model.Message{Type: in.Type, ID: in.Id, Data: data}
		cb(msg)
	}
}

// ClientID returns the client's id
func (c *GRPCRealtimeClient) ClientID() string {
	return c.id
}

// Context returns the client's context
func (c *GRPCRealtimeClient) Context() context.Context {
	return c.ctx
}
