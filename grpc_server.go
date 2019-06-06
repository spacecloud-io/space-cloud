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

	// Create a create request
	req := model.CreateRequest{}

	var temp interface{}
	if err := json.Unmarshal(in.Document, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Document = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := s.auth.IsCreateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the write operation
	err = s.crud.Create(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
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
				s.realtime.Send(&model.FeedData{
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
	return &pb.Response{Status: 200}, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.Response, error) {

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
	status, err := s.auth.IsReadOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the read operation
	result, err := s.crud.Read(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
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
	status, err := s.auth.IsUpdateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	err = s.crud.Update(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
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

			data, err := s.crud.Read(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &model.ReadRequest{Find: find, Operation: utils.One})
			if err == nil {
				s.realtime.Send(&model.FeedData{
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
	return &pb.Response{Status: 200}, nil

}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {

	// Load the request from the body
	req := model.DeleteRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Find = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := s.auth.IsDeleteOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the delete operation
	err = s.crud.Delete(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	// Send realtime message in dev mode
	if !s.isProd && req.Operation == utils.One {
		idVar := "id"
		if in.Meta.DbType == string(utils.Mongo) {
			idVar = "_id"
		}

		if id, p := req.Find[idVar]; p {
			s.realtime.Send(&model.FeedData{
				Group:     in.Meta.Col,
				Type:      utils.RealtimeDelete,
				TimeStamp: time.Now().Unix(),
				DocID:     id.(string),
				DBType:    in.Meta.DbType,
			})
		}
	}

	// Give positive acknowledgement
	return &pb.Response{Status: 200}, nil
}

func (s *server) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.Response, error) {

	req := model.AggregateRequest{}
	temp := []map[string]interface{}{}
	if err := json.Unmarshal(in.Pipeline, &temp); err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}
	req.Pipeline = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := s.auth.IsAggregateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the read operation
	result, err := s.crud.Aggregate(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
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

	allRequests := []model.AllRequest{}
	for _, req := range in.Batchrequest {
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
			status, err := s.auth.IsCreateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &r)
			if err != nil {
				return &pb.Response{Status: int32(status), Error: err.Error()}, nil
			}
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
			status, err := s.auth.IsUpdateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &r)
			if err != nil {
				return &pb.Response{Status: int32(status), Error: err.Error()}, nil
			}

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
			status, err := s.auth.IsDeleteOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &r)
			if err != nil {
				return &pb.Response{Status: int32(status), Error: err.Error()}, nil
			}
		}
	}
	// Perform the Batch operation
	batch := model.BatchRequest{}
	batch.Requests = allRequests
	err := s.crud.Batch(ctx, in.Meta.DbType, in.Meta.Project, &batch)
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
						s.realtime.Send(&model.FeedData{
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
							s.realtime.Send(&model.FeedData{
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

						data, err := s.crud.Read(ctx, in.Meta.DbType, in.Meta.Project, req.Col, &model.ReadRequest{Find: find, Operation: utils.One})
						if err == nil {
							s.realtime.Send(&model.FeedData{
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
	var params interface{}
	if err := json.Unmarshal(in.Params, &params); err != nil {
		out := pb.Response{}
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	auth, err := s.auth.IsFuncCallAuthorised(in.Project, in.Service, in.Function, in.Token, params)
	if err != nil {
		return &pb.Response{Status: 403, Error: err.Error()}, nil
	}

	result, err := s.functions.Call(in.Service, in.Function, auth, params, int(in.Timeout))
	if err != nil {
		return &pb.Response{Status: 500, Error: err.Error()}, nil
	}

	data, _ := json.Marshal(result)
	return &pb.Response{Result: data, Status: 200}, nil
}

func (s *server) Service(stream pb.SpaceCloud_ServiceServer) error {
	client := client.CreateGRPCServiceClient(stream)
	defer s.functions.UnregisterService(client.ClientID())
	defer client.Close()
	go client.RoutineWrite()

	// Get client details
	clientID := client.ClientID()

	client.Read(func(req *model.Message) {
		switch req.Type {
		case utils.TypeServiceRegister:
			// TODO add security rule for functions registered as well
			data := new(model.ServiceRegisterRequest)
			mapstructure.Decode(req.Data, data)

			s.functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
				client.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
			})

			client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

		case utils.TypeServiceRequest:
			data := new(model.FunctionsPayload)
			mapstructure.Decode(req.Data, data)

			s.functions.HandleServiceResponse(data)

		}
	})
	return nil
}

func (s *server) RealTime(stream pb.SpaceCloud_RealTimeServer) error {
	client := client.CreateGRPCRealtimeClient(stream)
	defer s.realtime.RemoveClient(client.ClientID())
	defer client.Close()
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

			// Subscribe to relaitme feed
			feedData, err := s.realtime.Subscribe(ctx, clientID, s.auth, s.crud, data, func(feed *model.FeedData) {
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

			s.realtime.Unsubscribe(clientID, data)

			// Send response to client
			res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
			client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
		}
	})
	return nil
}

func (s *server) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.Response, error) {
	status, result, err := s.user.Profile(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project, in.Id)
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
	status, result, err := s.user.Profiles(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project)
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
	status, result, err := s.user.EmailEditProfile(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project, in.Id, in.Email, in.Name, in.Password)
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
	status, result, err := s.user.EmailSignIn(ctx, in.Meta.DbType, in.Meta.Project, in.Email, in.Password)
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
	status, result, err := s.user.EmailSignUp(ctx, in.Meta.DbType, in.Meta.Project, in.Email, in.Name, in.Password, in.Role)
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
