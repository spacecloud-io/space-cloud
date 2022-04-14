package admin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (a *App) getGlobalTypes() model.OperationTypes {
	reflector := utils.GetJSONSchemaReflector()

	return model.OperationTypes{
		"env": &model.OperationTypeDefinition{
			Method:          http.MethodGet,
			RequiredParents: []string{},
			ResponseSchema:  reflector.Reflect(&loadEnvResponse{}),
			Controller: model.OperationHooks{
				Handle: func(ctx context.Context, obj *model.ResourceObject, reqParams *model.RequestParams) (int, interface{}, error) {
					return http.StatusOK, loadEnvResponse{IsProd: !a.IsDev, Version: utils.BuildVersion, ClusterType: "istio", LoginURL: "/mission-control/login"}, nil
				},
			},
		},
		"refresh-token": &model.OperationTypeDefinition{
			Method:          http.MethodPost,
			RequiredParents: []string{},
			RequestSchema:   reflector.Reflect(&tokenPayload{}),
			ResponseSchema:  reflector.Reflect(&tokenPayload{}),
			Controller: model.OperationHooks{
				DecodePayload: func(ctx context.Context, reader io.ReadCloser) (interface{}, error) {
					var req tokenPayload
					err := json.NewDecoder(reader).Decode(&req)
					return req, err
				},
				Handle: func(ctx context.Context, obj *model.ResourceObject, reqParams *model.RequestParams) (int, interface{}, error) {
					req := obj.Spec.(tokenPayload)

					// parse the token to extract the claims
					claims, err := a.auth.Verify(req.Token)
					if err != nil {
						return http.StatusBadRequest, model.ErrorResponse{Error: "Unable to verify admin token"}, nil
					}

					// Create a new token
					newToken, err := a.auth.Sign(claims)
					if err != nil {
						return http.StatusBadRequest, model.ErrorResponse{Error: "Unable to sign new admin token"}, nil
					}

					return http.StatusOK, tokenPayload{Token: newToken}, nil
				},
			},
		},
		"login": &model.OperationTypeDefinition{
			Method:          http.MethodPost,
			RequiredParents: []string{},
			RequestSchema:   reflector.Reflect(&loginRequest{}),
			ResponseSchema:  reflector.Reflect(&tokenPayload{}),
			Controller: model.OperationHooks{
				DecodePayload: func(ctx context.Context, reader io.ReadCloser) (interface{}, error) {
					var req loginRequest
					err := json.NewDecoder(reader).Decode(&req)
					return req, err
				},
				Handle: func(ctx context.Context, obj *model.ResourceObject, reqParams *model.RequestParams) (int, interface{}, error) {
					req := obj.Spec.(loginRequest)

					// Check if provided creds are valid
					if req.User != a.User || req.Key != a.Pass {
						return http.StatusNotFound, model.ErrorResponse{Error: "Invalid admin credentials provided"}, nil
					}

					// Generate a token
					token, err := a.auth.Sign(map[string]interface{}{"id": req.User, "role": "admin"})
					if err != nil {
						return 0, nil, err
					}

					// Return response
					return http.StatusOK, tokenPayload{Token: token}, nil
				},
			},
		},
	}
}
