package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
)

func getProjectConfigTypes() model.ConfigTypes {
	reflector := utils.GetJSONSchemaReflector()
	return model.ConfigTypes{
		"project": &model.ConfigTypeDefinition{
			Schema:          reflector.Reflect(&projectConfig{}),
			RequiredParents: []string{},
		},
		"aes-key": &model.ConfigTypeDefinition{
			IsSecure:        true,
			Schema:          reflector.Reflect(&projectAesKey{}),
			RequiredParents: []string{"project"},
			Controller: model.ConfigHooks{
				PreApply: func(ctx context.Context, obj *model.ResourceObject, store model.StoreMan) error {
					// Get all the aes keys in the project
					l, err := store.GetResources(ctx, &obj.Meta)
					if err != nil {
						return err
					}

					// Check if an aes key already exists in this project. We want to allow just
					// one aes key per project
					if len(l.List) > 0 && l.List[0].Meta.Name == obj.Meta.Name {
						return errors.New("only one aes key allowed per project")
					}
					return nil
				},
			},
		},
		"jwt-secret": &model.ConfigTypeDefinition{
			IsSecure:        true,
			Schema:          reflector.Reflect(&config.Secret{}),
			RequiredParents: []string{"project"},
		},
	}
}

func (a *App) getProjectOperationTypes() model.OperationTypes {
	reflector := utils.GetJSONSchemaReflector()
	return model.OperationTypes{
		"generate-internal-token": &model.OperationTypeDefinition{
			Method:          http.MethodGet,
			RequiredParents: []string{"project"},
			ResponseSchema:  reflector.Reflect(&tokenPayload{}),
			Controller: model.OperationHooks{
				Handle: func(ctx context.Context, obj *model.ResourceObject, reqParams *model.RequestParams) (int, interface{}, error) {
					// Check if the project exists
					projectID := obj.Meta.Parents["project"]
					project, p := a.Projects[projectID]
					if !p {
						return http.StatusBadRequest, model.ErrorResponse{Error: fmt.Sprintf("Provided project '%s' does not exist", projectID)}, nil
					}

					// Generate a internal token for the project
					token, err := project.auth.Sign(map[string]interface{}{"id": utils.InternalUserID, "claims": reqParams.Claims})
					if err != nil {
						return http.StatusInternalServerError, nil, err
					}

					return http.StatusOK, tokenPayload{Token: token}, nil
				},
			},
		},
	}
}
