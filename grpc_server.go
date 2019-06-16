package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/space-cloud/model"
	pb "github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/client"
)

func (s *server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	// Create a create request
	req := model.CreateRequest{}

	var temp interface{}
	if err := json.Unmarshal(in.Document, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Document = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := state.Auth.IsCreateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Send realtime message intent
	msgID := state.Realtime.SendCreateIntent(in.Meta.Project, in.Meta.DbType, in.Meta.Col, &req)

	// Perform the write operation
	err = state.Crud.Create(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		// Send realtime nack
		state.Realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, false)

		// Send gRPC Response
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Send realtime ack
	state.Realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, true)

	// Give positive acknowledgement
	return &pb.Response{Status: 200}, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	req := model.ReadRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Find = temp
	req.Operation = in.Operation
	opts := model.ReadOptions{}
	opts.Select = in.Options.Select
	opts.Sort = in.Options.Sort
	opts.Skip = &in.Options.Skip
	opts.Limit = &in.Options.Limit
	opts.Distinct = &in.Options.Distinct
	req.Options = &opts

	// Create empty read options if it does not exist
	if req.Options == nil {
		req.Options = new(model.ReadOptions)
	}

	// Check if the user is authenticated
	status, err := state.Auth.IsReadOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the read operation
	result, err := state.Crud.Read(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Give positive acknowledgement
	return &pb.Response{Status: 200, Result: resultBytes}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	req := model.UpdateRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Find = temp

	temp = map[string]interface{}{}
	if err := json.Unmarshal(in.Update, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Update = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := state.Auth.IsUpdateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Send realtime message intent
	msgID := state.Realtime.SendUpdateIntent(in.Meta.Project, in.Meta.DbType, in.Meta.Col, &req)

	err = state.Crud.Update(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		// Send realtime nack
		state.Realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, false)

		// Send gRPC Response
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Send realtime ack
	state.Realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, true)

	// Give positive acknowledgement
	return &pb.Response{Status: 200}, nil

}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	// Load the request from the body
	req := model.DeleteRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Find = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := state.Auth.IsDeleteOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Send realtime message intent
	msgID := state.Realtime.SendDeleteIntent(in.Meta.Project, in.Meta.DbType, in.Meta.Col, &req)

	// Perform the delete operation
	err = state.Crud.Delete(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		// Send realtime nack
		state.Realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, false)

		// Send gRPC Response
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Send realtime ack
	state.Realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, true)

	// Give positive acknowledgement
	return &pb.Response{Status: 200}, nil
}

func (s *server) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	req := model.AggregateRequest{}
	temp := []map[string]interface{}{}
	if err := json.Unmarshal(in.Pipeline, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Pipeline = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := state.Auth.IsAggregateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the read operation
	result, err := state.Crud.Aggregate(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Give positive acknowledgement
	return &pb.Response{Status: 200, Result: resultBytes}, nil
}

func (s *server) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	type msg struct {
		id, col string
	}

	msgIDs := make([]*msg, len(in.Batchrequest))

	allRequests := []model.AllRequest{}
	for i, req := range in.Batchrequest {
		switch req.Type {
		case string(utils.Create):
			eachReq := model.AllRequest{}
			eachReq.Type = req.Type
			eachReq.Col = req.Col

			r := model.CreateRequest{}
			var temp interface{}
			if err := json.Unmarshal(req.Document, &temp); err != nil {
				return &pb.Response{Status: 500, Error: err.Error()}, nil
			}
			r.Document = temp
			eachReq.Document = temp

			r.Operation = req.Operation
			eachReq.Operation = req.Operation

			allRequests = append(allRequests, eachReq)

			// Check if the user is authenticated
			status, err := state.Auth.IsCreateOpAuthorised(in.Meta.Project, in.Meta.DbType, req.Col, in.Meta.Token, &r)
			if err != nil {
				return &pb.Response{Status: int32(status), Error: err.Error()}, nil
			}

			// Send realtime message intent
			msgID := state.Realtime.SendCreateIntent(in.Meta.Project, in.Meta.DbType, req.Col, &r)
			msgIDs[i] = &msg{id: msgID, col: req.Col}

		case string(utils.Update):
			eachReq := model.AllRequest{}
			eachReq.Type = req.Type
			eachReq.Col = req.Col

			r := model.UpdateRequest{}
			temp := map[string]interface{}{}
			if err := json.Unmarshal(req.Find, &temp); err != nil {
				return &pb.Response{Status: 500, Error: err.Error()}, nil
			}
			r.Find = temp
			eachReq.Find = temp

			temp = map[string]interface{}{}
			if err := json.Unmarshal(req.Update, &temp); err != nil {
				return &pb.Response{Status: 500, Error: err.Error()}, nil
			}
			r.Update = temp
			eachReq.Update = temp

			r.Operation = req.Operation
			eachReq.Operation = req.Operation

			allRequests = append(allRequests, eachReq)

			// Check if the user is authenticated
			status, err := state.Auth.IsUpdateOpAuthorised(in.Meta.Project, in.Meta.DbType, req.Col, in.Meta.Token, &r)
			if err != nil {
				return &pb.Response{Status: int32(status), Error: err.Error()}, nil
			}

			// Send realtime message intent
			msgID := state.Realtime.SendUpdateIntent(in.Meta.Project, in.Meta.DbType, req.Col, &r)
			msgIDs[i] = &msg{id: msgID, col: req.Col}

		case string(utils.Delete):
			eachReq := model.AllRequest{}
			eachReq.Type = req.Type
			eachReq.Col = req.Col

			r := model.DeleteRequest{}
			temp := map[string]interface{}{}
			if err := json.Unmarshal(req.Find, &temp); err != nil {
				return &pb.Response{Status: 500, Error: err.Error()}, nil
			}
			r.Find = temp
			eachReq.Find = temp

			r.Operation = req.Operation
			eachReq.Operation = req.Operation

			allRequests = append(allRequests, eachReq)

			// Check if the user is authenticated
			status, err := state.Auth.IsDeleteOpAuthorised(in.Meta.Project, in.Meta.DbType, req.Col, in.Meta.Token, &r)
			if err != nil {
				return &pb.Response{Status: int32(status), Error: err.Error()}, nil
			}

			// Send realtime message intent
			msgID := state.Realtime.SendDeleteIntent(in.Meta.Project, in.Meta.DbType, req.Col, &r)
			msgIDs[i] = &msg{id: msgID, col: req.Col}
		}
	}
	// Perform the Batch operation
	batch := model.BatchRequest{}
	batch.Requests = allRequests
	err = state.Crud.Batch(ctx, in.Meta.DbType, in.Meta.Project, &batch)
	if err != nil {
		// Send realtime nack
		for _, m := range msgIDs {
			state.Realtime.SendAck(m.id, in.Meta.Project, m.col, false)
		}

		// Send gRPC Response
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Send realtime nack
	for _, m := range msgIDs {
		state.Realtime.SendAck(m.id, in.Meta.Project, m.col, true)
	}

	// Give positive acknowledgement
	return &pb.Response{Status: 200}, nil
}

func (s *server) Call(ctx context.Context, in *pb.FunctionsRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	var params interface{}
	if err := json.Unmarshal(in.Params, &params); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	auth, err := state.Auth.IsFuncCallAuthorised(in.Project, in.Service, in.Function, in.Token, params)
	if err != nil {
		return &pb.Response{Status: 403, Error: err.Error()}, nil
	}

	result, err := state.Functions.Call(in.Service, in.Function, auth, params, int(in.Timeout))
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	data, _ := json.Marshal(result)
	return &pb.Response{Result: data, Status: 200}, nil
}

func (s *server) Service(stream pb.SpaceCloud_ServiceServer) error {
	// Create an empty project variable
	var project string

	// Create a new client
	client := client.CreateGRPCServiceClient(stream)

	defer func() {
		// Unregister service if project could be loaded
		state, err := s.projects.LoadProject(project)
		if err == nil {
			// Unregister the service
			state.Functions.UnregisterService(client.ClientID())
		}
	}()

	// Close the client to free up resources
	defer client.Close()

	// Start the writer routine
	go client.RoutineWrite()

	// Get client details
	clientID := client.ClientID()

	client.Read(func(req *model.Message) {
		switch req.Type {
		case utils.TypeServiceRegister:
			// TODO add security rule for functions registered as well
			data := new(model.ServiceRegisterRequest)
			mapstructure.Decode(req.Data, data)

			// Set the clients project
			project = data.Project

			state, err := s.projects.LoadProject(project)
			if err != nil {
				client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": false}})
				return
			}
			state.Functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
				client.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
			})

			client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

		case utils.TypeServiceRequest:
			data := new(model.FunctionsPayload)
			mapstructure.Decode(req.Data, data)

			// Handle response if project could be loaded
			state, err := s.projects.LoadProject(project)
			if err == nil {
				state.Functions.HandleServiceResponse(data)
			}
		}
	})
	return nil
}

func (s *server) RealTime(stream pb.SpaceCloud_RealTimeServer) error {
	// Create an empty project variable
	var project string

	// Create a new client
	client := client.CreateGRPCRealtimeClient(stream)

	defer func() {
		// Unregister service if project could be loaded
		state, err := s.projects.LoadProject(project)
		if err == nil {
			// Unregister the service
			state.Realtime.RemoveClient(client.ClientID())
		}
	}()

	// Close the client to free up resources
	defer client.Close()

	// Start the writer routine
	go client.RoutineWrite()

	// Get client details
	ctx := client.Context()
	clientID := client.ClientID()

	client.Read(func(req *model.Message) {
		switch req.Type {
		case utils.TypeRealtimeSubscribe:

			// For realtime subscribe event
			data := new(model.RealtimeRequest)
			mapstructure.Decode(req.Data, data)

			// Set the clients project
			project = data.Project

			// Load the project state
			state, err := s.projects.LoadProject(project)
			if err != nil {
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
				client.Write(&model.Message{ID: req.ID, Type: utils.TypeRealtimeSubscribe, Data: res})
				return
			}

			// Subscribe to relaitme feed
			feedData, err := state.Realtime.Subscribe(ctx, clientID, state.Auth, state.Crud, data, func(feed *model.FeedData) {
				client.Write(&model.Message{ID: req.ID, Type: utils.TypeRealtimeFeed, Data: feed})
			})
			if err != nil {
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
				client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
				return
			}

			// Send response to client
			res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData}
			client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

		case utils.TypeRealtimeUnsubscribe:
			// For realtime subscribe event
			data := new(model.RealtimeRequest)
			mapstructure.Decode(req.Data, data)

			// Load the project state
			state, err := s.projects.LoadProject(project)
			if err != nil {
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
				client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
				return
			}

			state.Realtime.Unsubscribe(clientID, data)

			// Send response to client
			res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
			client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
		}
	})
	return nil
}

func (s *server) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	status, result, err := state.UserManagement.Profile(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project, in.Id)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	res, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = http.StatusInternalServerError
		out.Error = err1.Error()
		return &out, nil
	}
	out.Result = res
	return &out, nil
}

func (s *server) Profiles(ctx context.Context, in *pb.ProfilesRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	status, result, err := state.UserManagement.Profiles(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	res, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = http.StatusInternalServerError
		out.Error = err1.Error()
		return &out, nil
	}
	out.Result = res
	return &out, nil
}

func (s *server) EditProfile(ctx context.Context, in *pb.EditProfileRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	status, result, err := state.UserManagement.EmailEditProfile(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project, in.Id, in.Email, in.Name, in.Password)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	res, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = http.StatusInternalServerError
		out.Error = err1.Error()
		return &out, nil
	}
	out.Result = res
	return &out, nil
}

func (s *server) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	status, result, err := state.UserManagement.EmailSignIn(ctx, in.Meta.DbType, in.Meta.Project, in.Email, in.Password)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	res, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = http.StatusInternalServerError
		out.Error = err1.Error()
		return &out, nil
	}
	out.Result = res
	return &out, nil
}

func (s *server) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	status, result, err := state.UserManagement.EmailSignUp(ctx, in.Meta.DbType, in.Meta.Project, in.Email, in.Name, in.Password, in.Role)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	res, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = http.StatusInternalServerError
		out.Error = err1.Error()
		return &out, nil
	}
	out.Result = res
	return &out, nil
}
