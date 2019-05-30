package utils

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	pb "github.com/spaceuptech/space-cloud/proto"
)

// Client is the object handling all client interactions
type Client struct {
	id           string
	channel      chan *model.Message
	Context      context.Context
	cancel       context.CancelFunc
	socket       *websocket.Conn              //only for Websocket
	stream       pb.SpaceCloud_RealTimeServer //Only for grpc
	streamServer pb.SpaceCloud_ServiceServer  // Only for grpc (service)
	protocol     RealTimeProtocol
}

// DataCallback is the callback invoked when data is read by the socket
type DataCallback func(data *model.Message)

// RoutineWrite starts a json writer routine
func (c *Client) RoutineWrite() {
	switch c.protocol {
	case Websocket:
		for res := range c.channel {
			err := c.socket.WriteJSON(res)
			if err != nil {
				log.Println(err)
				return
			}
		}
	case GRPC:
		for res := range c.channel {
			//Convert the Message into RealTime response.
			switch res.Type {
			case TypeRealtimeSubscribe, TypeRealtimeUnsubscribe:
				//Decode the Message
				responseMsg := res.Data.(model.RealtimeResponse)
				feedData := make([]*pb.FeedData, len(responseMsg.Docs))
				for i, feed := range responseMsg.Docs {
					payload, err := json.Marshal(feed.Payload)
					if err != nil {
						log.Println(err)
						return
					}
					feedData[i] = &pb.FeedData{QueryId: feed.QueryID, DocId: feed.DocID, Type: feed.Type, Group: feed.Group, DbType: feed.DBType, Payload: payload, TimeStamp: feed.TimeStamp}
				}
				grpcResponse := pb.RealTimeResponse{Id: res.ID, Group: responseMsg.Group, Ack: responseMsg.Ack, Error: responseMsg.Error, FeedData: feedData}
				c.stream.Send(&grpcResponse)

			case TypeRealtimeFeed:
				feed := res.Data.(model.FeedData)
				feedData := make([]*pb.FeedData, 1)
				payload, err := json.Marshal(feed.Payload)
				if err != nil {
					log.Println(err)
					return
				}
				feedData[0] = &pb.FeedData{
					QueryId: feed.QueryID, DocId: feed.DocID, Type: feed.Type, Group: feed.Group, DbType: feed.DBType, Payload: payload, TimeStamp: feed.TimeStamp}
				grpcResponse := pb.RealTimeResponse{Id: res.ID, Group: res.Data.(model.FeedData).Group, FeedData: feedData}
				c.stream.Send(&grpcResponse)
			}
		}

	case GRPCService:
		for res := range c.channel {
			switch res.Type {
			case TypeServiceRequest:
				reqMsg, ok := res.Data.(*model.FunctionsPayload)
				if !ok {
					log.Println("GRPC Service Error - Invalid data type", res.Data)
					break
				}

				authData, _ := json.Marshal(reqMsg.Auth)
				paramsData, _ := json.Marshal(reqMsg.Params)
				c.streamServer.Send(&pb.FunctionsPayload{
					Auth:     authData,
					Params:   paramsData,
					Function: reqMsg.Func,
				})

			default:
				log.Println("GRPC Service Error - Invalid request type", res.Type)
			}
		}
	}
}

// Write wrties the object to the client
func (c *Client) Write(res *model.Message) {
	select {
	case c.channel <- res:
	case <-c.Context.Done():
	}
}

// Close closes the client
func (c *Client) Close() {
	c.cancel()
	close(c.channel)
	c.socket.Close()
}

// Read startes a blocking reader routine
func (c *Client) Read(cb DataCallback) {
	switch c.protocol {
	case Websocket:
		for {
			data := &model.Message{}
			err := c.socket.ReadJSON(data)
			if err != nil {
				log.Println(err)
				return
			}

			cb(data)
		}

	case GRPC:
		for {
			in, err := c.stream.Recv()
			if err != nil {
				log.Println("GRPC Error -", err)
				return
			}
			var data map[string]interface{}
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

	case GRPCService:
		for {
			in, err := c.streamServer.Recv()
			if err != nil {
				if err != nil {
					log.Println("GRPC Service Error -", err)
					return
				}
			}

			switch in.Type {
			case TypeServiceRegister:
				data := map[string]interface{}{"service": in.Service}
				msg := &model.Message{ID: in.Id, Type: TypeServiceRegister, Data: data}
				cb(msg)

			case TypeServiceRequest:
				var params interface{}
				json.Unmarshal(in.Params, &params)
				data := map[string]interface{}{
					"params": params,
				}
				msg := &model.Message{ID: in.Id, Type: TypeServiceRequest, Data: data}
				cb(msg)

			default:
				log.Println("GRPC Service Error - Invalid request type", in.Type)
			}
		}
	}
}

// ClientID returns the client's id
func (c *Client) ClientID() string {
	return c.id
}

// CreateWebsocketClient makes a client object to manage the socket
func CreateWebsocketClient(socket *websocket.Conn) *Client {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &Client{id, channel, ctx, cancel, socket, nil, nil, Websocket}
}

// CreateGRPCClient makes a client object to manage the grpc
func CreateGRPCClient(stream pb.SpaceCloud_RealTimeServer) *Client {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &Client{id, channel, ctx, cancel, nil, stream, nil, GRPC}
}

// CreateGRPCServiceClient makes a client object to manage the grpc
func CreateGRPCServiceClient(stream pb.SpaceCloud_ServiceServer) *Client {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &Client{id, channel, ctx, cancel, nil, nil, stream, GRPCService}
}
