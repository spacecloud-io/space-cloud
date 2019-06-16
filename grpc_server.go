package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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

	authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, in.Meta.Col, utils.Create)
	if err != nil {
		out := pb.Response{}
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.CreateRequest{}
	if in.Operation == utils.One {
		temp := map[string]interface{}{}
		if err = json.Unmarshal(in.Document, &temp); err != nil {
			out := pb.Response{}
			out.Status = 500
			out.Error = err.Error()
			return &out, nil
		}
		req.Document = temp
	} else if in.Operation == utils.All {
		temp := []interface{}{}
		if err = json.Unmarshal(in.Document, &temp); err != nil {
			out := pb.Response{}
			out.Status = 500
			out.Error = err.Error()
			return &out, nil
		}
		req.Document = temp
	}
	req.Operation = in.Operation

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"doc": req.Document, "op": req.Operation, "auth": authObj},
		"project": in.Meta.Project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, in.Meta.Col, utils.Create, args)
	if err != nil {
		out := pb.Response{}
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the write operation
	err = state.Crud.Create(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	// Send realtime message in dev mode
	if !s.isProd {
		var rows []interface{}
		switch req.Operation {
		case utils.One:
			rows = []interface{}{req.Document}
		case utils.All:
			rows = req.Document.([]interface{})
		default:
			rows = []interface{}{}
		}

		for _, t := range rows {
			data := t.(map[string]interface{})

			idVar := "id"
			if in.Meta.DbType == string(utils.Mongo) {
				idVar = "_id"
			}

			// Send realtime message if id fields exists
			if id, p := data[idVar]; p {
				state.Realtime.Send(&model.FeedData{
					Group:     in.Meta.Col,
					DBType:    in.Meta.DbType,
					Type:      utils.RealtimeWrite,
					TimeStamp: time.Now().Unix(),
					DocID:     id.(string),
					Payload:   data,
				})
			}
		}
	}

	// Give positive acknowledgement
	out := pb.Response{}
	out.Status = 200
	return &out, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, in.Meta.Col, utils.Read)
	if err != nil {
		out := pb.Response{}
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.ReadRequest{}
	temp := map[string]interface{}{}
	if err = json.Unmarshal(in.Find, &temp); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
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

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
		"project": in.Meta.Project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, in.Meta.Col, utils.Read, args)
	if err != nil {
		out := pb.Response{}
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the read operation
	result, err := state.Crud.Read(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	temp1, err1 := json.Marshal(result)
	if err1 != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err1.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out := pb.Response{}
	out.Status = 200
	out.Result = temp1
	return &out, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	// Check if the user is authicated
	authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, in.Meta.Col, utils.Update)
	if err != nil {
		out := pb.Response{}
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.UpdateRequest{}
	temp := map[string]interface{}{}
	if err = json.Unmarshal(in.Find, &temp); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}
	req.Find = temp

	temp = map[string]interface{}{}
	if err = json.Unmarshal(in.Update, &temp); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}
	req.Update = temp
	req.Operation = in.Operation

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
		"project": in.Meta.Project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, in.Meta.Col, utils.Read, args)
	if err != nil {
		out := pb.Response{}
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	err = state.Crud.Update(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	// Send realtime message in dev mode
	if !s.isProd && req.Operation == utils.One {
		idVar := "id"
		if in.Meta.DbType == string(utils.Mongo) {
			idVar = "_id"
		}

		if id, p := req.Find[idVar]; p {
			// Create the find object
			find := map[string]interface{}{}

			switch utils.DBType(in.Meta.DbType) {
			case utils.Mongo:
				find["_id"] = id

			default:
				find["id"] = id
			}

			data, err := state.Crud.Read(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &model.ReadRequest{Find: find, Operation: utils.One})
			if err == nil {
				state.Realtime.Send(&model.FeedData{
					Group:     in.Meta.Col,
					Type:      utils.RealtimeWrite,
					TimeStamp: time.Now().Unix(),
					DocID:     id.(string),
					DBType:    in.Meta.DbType,
					Payload:   data.(map[string]interface{}),
				})
			}
		}
	}

	// Give positive acknowledgement
	out := pb.Response{}
	out.Status = 200
	return &out, nil

}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	// Check if the user is authicated
	authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, in.Meta.Col, utils.Delete)
	if err != nil {
		out := pb.Response{}
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	// Load the request from the body
	req := model.DeleteRequest{}
	temp := map[string]interface{}{}
	if err = json.Unmarshal(in.Find, &temp); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}
	req.Find = temp
	req.Operation = in.Operation

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": req.Find, "op": req.Operation, "auth": authObj},
		"project": in.Meta.Project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, in.Meta.Col, utils.Delete, args)
	if err != nil {
		out := pb.Response{}
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the delete operation
	err = state.Crud.Delete(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	// Send realtime message in dev mode
	if !s.isProd && req.Operation == utils.One {
		idVar := "id"
		if in.Meta.DbType == string(utils.Mongo) {
			idVar = "_id"
		}

		if id, p := req.Find[idVar]; p {
			state.Realtime.Send(&model.FeedData{
				Group:     in.Meta.Col,
				Type:      utils.RealtimeDelete,
				TimeStamp: time.Now().Unix(),
				DocID:     id.(string),
				DBType:    in.Meta.DbType,
			})
		}
	}

	// Give positive acknowledgement
	out := pb.Response{}
	out.Status = 200
	return &out, nil
}

func (s *server) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	// Check if the user is authicated
	authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, in.Meta.Col, utils.Delete)
	if err != nil {
		out := pb.Response{}
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.AggregateRequest{}
	temp := []map[string]interface{}{}
	if err = json.Unmarshal(in.Pipeline, &temp); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}
	req.Pipeline = temp
	req.Operation = in.Operation

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": req.Pipeline, "op": req.Operation, "auth": authObj},
		"project": in.Meta.Project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, in.Meta.Col, utils.Aggregation, args)
	if err != nil {
		out := pb.Response{}
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the read operation
	result, err := state.Crud.Aggregate(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	temp1, err1 := json.Marshal(result)
	if err1 != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err1.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out := pb.Response{}
	out.Status = 200
	out.Result = temp1
	return &out, nil
}

func (s *server) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.Response, error) {
	// Load the project state
	state, err := s.projects.LoadProject(in.Meta.Project)
	if err != nil {
		return &pb.Response{Status: 400, Error: err.Error()}, nil
	}

	allRequests := []model.AllRequest{}
	for _, req := range in.Batchrequest {
		switch req.Type {

		case string(utils.Update):
			eachReq := model.AllRequest{}
			var updateObj map[string]interface{}
			if err := json.Unmarshal(req.Update, &updateObj); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			var findObj map[string]interface{}
			if err := json.Unmarshal(req.Update, &findObj); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			var document interface{}
			if err := json.Unmarshal(req.Update, &document); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			eachReq.Find = findObj
			eachReq.Update = updateObj
			eachReq.Document = document
			eachReq.Col = req.Col
			eachReq.Operation = req.Operation
			eachReq.Type = req.Type

			allRequests = append(allRequests, eachReq)

			authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, eachReq.Col, utils.Update)
			if err != nil {
				out := pb.Response{}
				out.Status = 401
				out.Error = err.Error()
				return &out, nil
			}
			args := map[string]interface{}{
				"args":    map[string]interface{}{"find": eachReq.Find, "update": eachReq.Update, "op": eachReq.Operation, "auth": authObj},
				"project": in.Meta.Project, // Don't forget to do this for every request
			}

			// Check if user is authorized to make this request
			err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, eachReq.Col, utils.Update, args)
			if err != nil {
				out := pb.Response{}
				out.Status = 403
				out.Error = err.Error()
				return &out, nil
			}

		case string(utils.Create):
			eachReq := model.AllRequest{}
			var updateObj map[string]interface{}
			if err := json.Unmarshal(req.Update, &updateObj); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			var findObj map[string]interface{}
			if err := json.Unmarshal(req.Update, &findObj); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			var document interface{}
			if err := json.Unmarshal(req.Update, &document); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			eachReq.Find = findObj
			eachReq.Update = updateObj
			eachReq.Document = document
			eachReq.Col = req.Col
			eachReq.Operation = req.Operation
			eachReq.Type = req.Type

			allRequests = append(allRequests, eachReq)

			authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, eachReq.Col, utils.Create)
			if err != nil {
				out := pb.Response{}
				out.Status = 401
				out.Error = err.Error()
				return &out, nil
			}
			// Create an args object
			args := map[string]interface{}{
				"args":    map[string]interface{}{"doc": eachReq.Document, "op": eachReq.Operation, "auth": authObj},
				"project": in.Meta.Project, // Don't forget to do this for every request
			}

			// Check if user is authorized to make this request
			err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, eachReq.Col, utils.Create, args)
			if err != nil {
				out := pb.Response{}
				out.Status = 403
				out.Error = err.Error()
				return &out, nil
			}

		case string(utils.Delete):
			eachReq := model.AllRequest{}
			var updateObj map[string]interface{}
			if err := json.Unmarshal(req.Update, &updateObj); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			var findObj map[string]interface{}
			if err := json.Unmarshal(req.Update, &findObj); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			var document interface{}
			if err := json.Unmarshal(req.Update, &document); err != nil {
				out := pb.Response{}
				out.Status = 500
				out.Error = err.Error()
				return &out, nil
			}
			eachReq.Find = findObj
			eachReq.Update = updateObj
			eachReq.Document = document
			eachReq.Col = req.Col
			eachReq.Operation = req.Operation
			eachReq.Type = req.Type

			allRequests = append(allRequests, eachReq)

			authObj, err := state.Auth.IsAuthenticated(in.Meta.Token, in.Meta.DbType, eachReq.Col, utils.Delete)
			if err != nil {
				out := pb.Response{}
				out.Status = 401
				out.Error = err.Error()
				return &out, nil
			}
			// Create an args object
			args := map[string]interface{}{
				"args":    map[string]interface{}{"find": eachReq.Find, "op": eachReq.Operation, "auth": authObj},
				"project": in.Meta.Project, // Don't forget to do this for every request
			}

			// Check if user is authorized to make this request
			err = state.Auth.IsAuthorized(in.Meta.Project, in.Meta.DbType, eachReq.Col, utils.Delete, args)
			if err != nil {
				out := pb.Response{}
				out.Status = 403
				out.Error = err.Error()
				return &out, nil
			}
		}
	}
	// Perform the Batch operation
	batch := model.BatchRequest{}
	batch.Requests = allRequests
	err = state.Crud.Batch(ctx, in.Meta.DbType, in.Meta.Project, &batch)
	if err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}
	if !s.isProd {

		for _, req := range batch.Requests {
			switch req.Type {
			case string(utils.Create):
				var rows []interface{}
				switch req.Operation {
				case utils.One:
					rows = []interface{}{req.Document}
				case utils.All:
					rows = req.Document.([]interface{})
				default:
					rows = []interface{}{}
				}

				for _, t := range rows {
					data := t.(map[string]interface{})

					idVar := "id"
					if in.Meta.DbType == string(utils.Mongo) {
						idVar = "_id"
					}

					// Send realtime message if id fields exists
					if id, p := data[idVar]; p {
						state.Realtime.Send(&model.FeedData{
							Group:     req.Col,
							DBType:    in.Meta.DbType,
							Type:      utils.RealtimeWrite,
							TimeStamp: time.Now().Unix(),
							DocID:     id.(string),
							Payload:   data,
						})
					}
				}

			case string(utils.Delete):
				if req.Operation == utils.One {
					idVar := "id"
					if in.Meta.DbType == string(utils.Mongo) {
						idVar = "_id"
					}

					if id, p := req.Find[idVar]; p {
						if err != nil {
							state.Realtime.Send(&model.FeedData{
								Group:     req.Col,
								Type:      utils.RealtimeDelete,
								TimeStamp: time.Now().Unix(),
								DocID:     id.(string),
								DBType:    in.Meta.DbType,
							})
						}
					}
				}

			case string(utils.Update):
				// Send realtime message in dev mode
				if req.Operation == utils.One {

					idVar := "id"
					if in.Meta.DbType == string(utils.Mongo) {
						idVar = "_id"
					}

					if id, p := req.Find[idVar]; p {
						// Create the find object
						find := map[string]interface{}{idVar: id}

						data, err := state.Crud.Read(ctx, in.Meta.DbType, in.Meta.Project, req.Col, &model.ReadRequest{Find: find, Operation: utils.One})
						if err == nil {
							state.Realtime.Send(&model.FeedData{
								Group:     req.Col,
								Type:      utils.RealtimeWrite,
								TimeStamp: time.Now().Unix(),
								DocID:     id.(string),
								DBType:    in.Meta.DbType,
								Payload:   data.(map[string]interface{}),
							})
						}
					}
				}
			}
		}
	}
	// Give positive acknowledgement
	out := pb.Response{}
	out.Status = 200
	return &out, nil
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
			// Subscribe to the realtime feed
			feedData, err := state.Realtime.Subscribe(ctx, clientID, state.Auth, state.Crud, data, func(feed *model.FeedData) {
				client.Write(&model.Message{Type: utils.TypeRealtimeFeed, Data: feed})
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
