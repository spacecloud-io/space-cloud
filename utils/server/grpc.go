package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	pb "github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/client"
)

// Create inserts document(s) into the database
func (s *Server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Response, error) {
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

	// Send realtime message intent
	msgID := s.realtime.SendCreateIntent(in.Meta.Project, in.Meta.DbType, in.Meta.Col, &req)

	// Perform the write operation
	err = s.crud.Create(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		// Send realtime nack
		s.realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, false)

		// Send gRPC Response
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	// Send realtime ack
	s.realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, true)

	// Give positive acknowledgement
	return &pb.Response{Status: http.StatusOK}, nil
}

// Read queries document(s) from the database
func (s *Server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.Response, error) {
	req := model.ReadRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
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
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	// Give positive acknowledgement
	return &pb.Response{Status: http.StatusOK, Result: resultBytes}, nil
}

// Update updates document(s) from the database
func (s *Server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.Response, error) {
	req := model.UpdateRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	req.Find = temp

	temp = map[string]interface{}{}
	if err := json.Unmarshal(in.Update, &temp); err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	req.Update = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := s.auth.IsUpdateOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Send realtime message intent
	msgID := s.realtime.SendUpdateIntent(in.Meta.Project, in.Meta.DbType, in.Meta.Col, &req)

	err = s.crud.Update(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		// Send realtime nack
		s.realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, false)

		// Send gRPC Response
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	// Send realtime ack
	s.realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, true)

	// Give positive acknowledgement
	return &pb.Response{Status: http.StatusOK}, nil
}

// Delete deletes document(s) from the database
func (s *Server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {
	// Load the request from the body
	req := model.DeleteRequest{}
	temp := map[string]interface{}{}
	if err := json.Unmarshal(in.Find, &temp); err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}
	req.Find = temp
	req.Operation = in.Operation

	// Check if the user is authenticated
	status, err := s.auth.IsDeleteOpAuthorised(in.Meta.Project, in.Meta.DbType, in.Meta.Col, in.Meta.Token, &req)
	if err != nil {
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Send realtime message intent
	msgID := s.realtime.SendDeleteIntent(in.Meta.Project, in.Meta.DbType, in.Meta.Col, &req)

	// Perform the delete operation
	err = s.crud.Delete(ctx, in.Meta.DbType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		// Send realtime nack
		s.realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, false)

		// Send gRPC Response
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	// Send realtime ack
	s.realtime.SendAck(msgID, in.Meta.Project, in.Meta.Col, true)

	// Give positive acknowledgement
	return &pb.Response{Status: http.StatusOK}, nil
}

// Aggregate aggregates document(s) from the database
func (s *Server) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.Response, error) {
	req := model.AggregateRequest{}
	temp := []map[string]interface{}{}
	if err := json.Unmarshal(in.Pipeline, &temp); err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
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
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	// Give positive acknowledgement
	return &pb.Response{Status: http.StatusOK, Result: resultBytes}, nil
}

// Batch performs a batch operation on the database
func (s *Server) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.Response, error) {
	type msg struct {
		id, col string
	}

	msgIDs := make([]*msg, len(in.Batchrequest))

	allRequests := []model.AllRequest{}
	for i, req := range in.Batchrequest {
		// Make status and error variables
		var status int
		var err error

		switch req.Type {
		case string(utils.Create):
			eachReq := model.AllRequest{}
			eachReq.Type = req.Type
			eachReq.Col = req.Col

			r := model.CreateRequest{}
			var temp interface{}
			if err = json.Unmarshal(req.Document, &temp); err != nil {
				status = http.StatusInternalServerError
			}
			r.Document = temp
			eachReq.Document = temp

			r.Operation = req.Operation
			eachReq.Operation = req.Operation

			allRequests = append(allRequests, eachReq)

			// Check if the user is authenticated
			status, err = s.auth.IsCreateOpAuthorised(in.Meta.Project, in.Meta.DbType, req.Col, in.Meta.Token, &r)
			if err == nil {
				// Send realtime message intent
				msgID := s.realtime.SendCreateIntent(in.Meta.Project, in.Meta.DbType, req.Col, &r)
				msgIDs[i] = &msg{id: msgID, col: req.Col}
			}

		case string(utils.Update):
			eachReq := model.AllRequest{}
			eachReq.Type = req.Type
			eachReq.Col = req.Col

			r := model.UpdateRequest{}
			temp := map[string]interface{}{}
			if err := json.Unmarshal(req.Find, &temp); err != nil {
				return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
			}
			r.Find = temp
			eachReq.Find = temp

			temp = map[string]interface{}{}
			if err = json.Unmarshal(req.Update, &temp); err != nil {
				status = http.StatusInternalServerError
			}
			r.Update = temp
			eachReq.Update = temp

			r.Operation = req.Operation
			eachReq.Operation = req.Operation

			allRequests = append(allRequests, eachReq)

			// Check if the user is authenticated
			status, err = s.auth.IsUpdateOpAuthorised(in.Meta.Project, in.Meta.DbType, req.Col, in.Meta.Token, &r)
			if err == nil {
				// Send realtime message intent
				msgID := s.realtime.SendUpdateIntent(in.Meta.Project, in.Meta.DbType, req.Col, &r)
				msgIDs[i] = &msg{id: msgID, col: req.Col}
			}

		case string(utils.Delete):
			eachReq := model.AllRequest{}
			eachReq.Type = req.Type
			eachReq.Col = req.Col

			r := model.DeleteRequest{}
			temp := map[string]interface{}{}
			if err = json.Unmarshal(req.Find, &temp); err != nil {
				status = http.StatusInternalServerError
			}
			r.Find = temp
			eachReq.Find = temp

			r.Operation = req.Operation
			eachReq.Operation = req.Operation

			allRequests = append(allRequests, eachReq)

			// Check if the user is authenticated
			status, err = s.auth.IsDeleteOpAuthorised(in.Meta.Project, in.Meta.DbType, req.Col, in.Meta.Token, &r)
			if err == nil {
				// Send realtime message intent
				msgID := s.realtime.SendDeleteIntent(in.Meta.Project, in.Meta.DbType, req.Col, &r)
				msgIDs[i] = &msg{id: msgID, col: req.Col}
			}
		}

		// Send negative acks and send error response
		for j := 0; j < i; j++ {
			s.realtime.SendAck(msgIDs[j].id, in.Meta.Project, msgIDs[j].col, false)
		}

		if err != nil {
			return &pb.Response{Status: int32(status), Error: err.Error()}, nil
		}

		// Send gRPC Response
		return &pb.Response{Status: int32(status), Error: err.Error()}, nil
	}

	// Perform the Batch operation
	batch := model.BatchRequest{}
	batch.Requests = allRequests
	err := s.crud.Batch(ctx, in.Meta.DbType, in.Meta.Project, &batch)
	if err != nil {
		// Send realtime nack
		for _, m := range msgIDs {
			s.realtime.SendAck(m.id, in.Meta.Project, m.col, false)
		}

		// Send gRPC Response
		return &pb.Response{Status: http.StatusInternalServerError, Error: err.Error()}, nil
	}

	// Send realtime nack
	for _, m := range msgIDs {
		s.realtime.SendAck(m.id, in.Meta.Project, m.col, true)
	}

	// Give positive acknowledgement
	return &pb.Response{Status: http.StatusOK}, nil
}

// Call invokes a function on the provided services
func (s *Server) Call(ctx context.Context, in *pb.FunctionsRequest) (*pb.Response, error) {
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

// Service registers and handles all opertions of a service
func (s *Server) Service(stream pb.SpaceCloud_ServiceServer) error {
	c := client.CreateGRPCServiceClient(stream)
	defer s.functions.UnregisterService(c.ClientID())
	defer c.Close()
	go c.RoutineWrite()

	// Get GRPC Service client details
	clientID := c.ClientID()

	c.Read(func(req *model.Message) bool {
		switch req.Type {
		case utils.TypeServiceRegister:
			// TODO add security rule for functions registered as well
			data := new(model.ServiceRegisterRequest)
			mapstructure.Decode(req.Data, data)

			s.functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
				c.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
			})

			c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

		case utils.TypeServiceRequest:
			data := new(model.FunctionsPayload)
			mapstructure.Decode(req.Data, data)

			s.functions.HandleServiceResponse(data)
		}

		return true
	})
	return nil
}

// RealTime registers and handles all opertions of a live query
func (s *Server) RealTime(stream pb.SpaceCloud_RealTimeServer) error {
	c := client.CreateGRPCRealtimeClient(stream)
	defer s.realtime.RemoveClient(c.ClientID())
	defer c.Close()
	go c.RoutineWrite()

	// Get GRPC Service client details
	ctx := c.Context()
	clientID := c.ClientID()

	c.Read(func(req *model.Message) bool {
		switch req.Type {
		case utils.TypeRealtimeSubscribe:
			// For realtime subscribe event
			data := new(model.RealtimeRequest)
			mapstructure.Decode(req.Data, data)

			// Subscribe to realtime feed
			feedData, err := s.realtime.Subscribe(ctx, clientID, s.auth, s.crud, data, func(feed *model.FeedData) {
				c.Write(&model.Message{ID: req.ID, Type: utils.TypeRealtimeFeed, Data: feed})
			})
			if err != nil {
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
				return true
			}

			// Send response to c
			res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData}
			c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

		case utils.TypeRealtimeUnsubscribe:
			// For realtime subscribe event
			data := new(model.RealtimeRequest)
			mapstructure.Decode(req.Data, data)

			s.realtime.Unsubscribe(clientID, data)

			// Send response to c
			res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
			c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
		}

		return true
	})

	return nil
}

// Profile queries the user's profiles
func (s *Server) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.Response, error) {
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

// Profiles queries all user profiles
func (s *Server) Profiles(ctx context.Context, in *pb.ProfilesRequest) (*pb.Response, error) {
	status, result, err := s.user.Profiles(ctx, in.Meta.Token, in.Meta.DbType, in.Meta.Project)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	res, err1 := json.Marshal(result["users"])
	if err1 != nil {
		out.Status = http.StatusInternalServerError
		out.Error = err1.Error()
		return &out, nil
	}
	out.Result = res

	return &out, nil
}

// EditProfile edits a user's profiles
func (s *Server) EditProfile(ctx context.Context, in *pb.EditProfileRequest) (*pb.Response, error) {
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

// SignIn signs a user in
func (s *Server) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.Response, error) {
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

// SignUp signs up a user
func (s *Server) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.Response, error) {
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

// CreateFolder creates a new folder
func (s *Server) CreateFolder(ctx context.Context, in *pb.CreateFolderRequest) (*pb.Response, error) {
	status, err := s.file.CreateDir(in.Meta.Project, in.Meta.Token, &model.CreateFileRequest{Name: in.Name, Path: in.Path, Type: "dir", MakeAll: false})
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	out.Result = []byte("")

	return &out, nil
}

// DeleteFile delete a file
func (s *Server) DeleteFile(ctx context.Context, in *pb.DeleteFileRequest) (*pb.Response, error) {
	status, err := s.file.DeleteFile(in.Meta.Project, in.Meta.Token, in.Path)
	out := pb.Response{}
	out.Status = int32(status)
	if err != nil {
		out.Error = err.Error()
		return &out, nil
	}
	out.Result = []byte("")

	return &out, nil
}

// ListFiles lists all files in the provided folder
func (s *Server) ListFiles(ctx context.Context, in *pb.ListFilesRequest) (*pb.Response, error) {
	status, result, err := s.file.ListFiles(in.Meta.Project, in.Meta.Token, &model.ListFilesRequest{Path: in.Path, Type: "all"})
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

// UploadFile uploads a file
func (s *Server) UploadFile(stream pb.SpaceCloud_UploadFileServer) error {
	req, err := stream.Recv()
	if err != nil {
		return stream.SendAndClose(&pb.Response{Status: int32(http.StatusInternalServerError), Error: err.Error()})
	}

	c := make(chan int)
	r, w := io.Pipe()
	// defer r.Close()
	// defer w.Close()

	go func() {
		status, err1 := s.file.UploadFile(req.Meta.Project, req.Meta.Token, &model.CreateFileRequest{Path: req.Path, Name: req.Name, Type: "file", MakeAll: true}, r)
		c <- status
		if err1 != nil {
			err = err1
		}
		w.Close()
	}()

	go func() {
		for {
			req, err1 := stream.Recv()
			if err1 == io.EOF {
				break
			}
			if err1 != nil {
				err = err1
				c <- http.StatusInternalServerError
				break
			}
			w.Write(req.Payload)
		}
		w.Close()
	}()

	status := <-c
	if err != nil {
		return stream.SendAndClose(&pb.Response{Status: int32(status), Error: err.Error()})
	}
	return stream.SendAndClose(&pb.Response{Status: int32(status), Result: []byte("")})
}

// DownloadFile downloads a file
func (s *Server) DownloadFile(in *pb.DownloadFileRequest, stream pb.SpaceCloud_DownloadFileServer) error {
	status, file, err := s.file.DownloadFile(in.Meta.Project, in.Meta.Token, in.Path)
	if err != nil {
		stream.Send(&pb.FilePayload{Status: int32(status), Error: err.Error()})
		return nil
	}
	defer file.Close()

	buf := make([]byte, utils.PayloadSize)
	for {
		n, err := file.File.Read(buf)
		if n > 0 {
			buf = buf[:n]
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			stream.Send(&pb.FilePayload{Status: int32(http.StatusInternalServerError), Error: err.Error()})
			break
		}
		req := pb.FilePayload{Payload: buf, Status: int32(http.StatusOK)}
		stream.Send(&req)
	}
	return nil
}
