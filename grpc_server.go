package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/spaceuptech/space-cloud/model"
	pb "github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
)

func (s *server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Response, error) {

	out := pb.Response{}

	authObj, err := s.auth.IsAuthenticated(in.Meta.Token, in.Meta.DBType, in.Meta.Col, utils.Create)
	if err != nil {
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.CreateRequest{}
	if in.Operation == utils.One {
		temp := map[string]interface{}{}
		if err = json.Unmarshal(in.Document, &temp); err != nil {
			out.Status = 500
			out.Error = err.Error()
			return &out, nil
		}
		req.Document = temp
	} else if in.Operation == utils.All {
		temp := []map[string]interface{}{}
		if err = json.Unmarshal(in.Document, &temp); err != nil {
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
	err = s.auth.IsAuthorized(in.Meta.DBType, in.Meta.Col, utils.Create, args)
	if err != nil {
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the write operation
	err = s.crud.Create(ctx, in.Meta.DBType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out.Status = 200
	return &out, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.Response, error) {

	out := pb.Response{}

	authObj, err := s.auth.IsAuthenticated(in.Meta.Token, in.Meta.DBType, in.Meta.Col, utils.Read)
	if err != nil {
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.ReadRequest{}
	temp := map[string]interface{}{}
	if err = json.Unmarshal(in.Find, &temp); err != nil {
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
	err = s.auth.IsAuthorized(in.Meta.DBType, in.Meta.Col, utils.Read, args)
	if err != nil {
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the read operation
	result, err := s.crud.Read(ctx, in.Meta.DBType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	temp1, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = 500
		out.Error = err1.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out.Status = 200
	out.Result = temp1
	return &out, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.Response, error) {

	out := pb.Response{}

	// Check if the user is authicated
	authObj, err := s.auth.IsAuthenticated(in.Meta.Token, in.Meta.DBType, in.Meta.Col, utils.Update)
	if err != nil {
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.UpdateRequest{}
	temp := map[string]interface{}{}
	if err = json.Unmarshal(in.Find, &temp); err != nil {
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}
	req.Find = temp
	if err = json.Unmarshal(in.Update, &temp); err != nil {
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
	err = s.auth.IsAuthorized(in.Meta.DBType, in.Meta.Col, utils.Read, args)
	if err != nil {
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	err = s.crud.Update(ctx, in.Meta.DBType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out.Status = 200
	return &out, nil

}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {

	out := pb.Response{}

	// Check if the user is authicated
	authObj, err := s.auth.IsAuthenticated(in.Meta.Token, in.Meta.DBType, in.Meta.Col, utils.Delete)
	if err != nil {
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	// Load the request from the body
	req := model.DeleteRequest{}
	temp := map[string]interface{}{}
	if err = json.Unmarshal(in.Find, &temp); err != nil {
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
	err = s.auth.IsAuthorized(in.Meta.DBType, in.Meta.Col, utils.Delete, args)
	if err != nil {
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the delete operation
	err = s.crud.Delete(ctx, in.Meta.DBType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out.Status = 200
	return &out, nil
}

func (s *server) Aggregate(ctx context.Context, in *pb.AggregateRequest) (*pb.Response, error) {

	out := pb.Response{}

	// Check if the user is authicated
	authObj, err := s.auth.IsAuthenticated(in.Meta.Token, in.Meta.DBType, in.Meta.Col, utils.Delete)
	if err != nil {
		out.Status = 401
		out.Error = err.Error()
		return &out, nil
	}

	req := model.AggregateRequest{}
	temp := []map[string]interface{}{}
	if err = json.Unmarshal(in.Pipeline, &temp); err != nil {
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
	err = s.auth.IsAuthorized(in.Meta.DBType, in.Meta.Col, utils.Aggregation, args)
	if err != nil {
		out.Status = 403
		out.Error = err.Error()
		return &out, nil
	}

	// Perform the read operation
	result, err := s.crud.Aggregate(ctx, in.Meta.DBType, in.Meta.Project, in.Meta.Col, &req)
	if err != nil {
		out.Status = 500
		out.Error = err.Error()
		return &out, nil
	}

	temp1, err1 := json.Marshal(result)
	if err1 != nil {
		out.Status = 500
		out.Error = err1.Error()
		return &out, nil
	}

	// Give positive acknowledgement
	out.Status = 200
	out.Result = temp1
	return &out, nil

}

//Rename the Function name.
func (s *server) initGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSpaceCloudServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
