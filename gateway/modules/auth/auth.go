package auth

import (
	"context"
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	jwtUtils "github.com/spaceuptech/space-cloud/gateway/utils/jwt"
)

// Module is responsible for authentication and authorisation
type Module struct {
	sync.RWMutex
	dbRules           config.DatabaseRules
	dbPrepQueryRules  config.DatabasePreparedQueries
	clusterID         string
	nodeID            string
	jwt               *jwtUtils.JWT
	crud              model.CrudAuthInterface
	fileRules         []*config.FileRule
	funcRules         config.Services
	securityFunctions config.SecurityFunctions
	eventingRules     map[string]*config.Rule
	project           string
	fileStoreType     string
	makeHTTPRequest   utils.TypeMakeHTTPRequest
	aesKey            []byte

	// Admin Manager
	adminMan adminMan
}

// Init creates a new instance of the auth object
func Init(clusterID, nodeID string, crud model.CrudAuthInterface, adminMan adminMan) *Module {
	return &Module{clusterID: clusterID, nodeID: nodeID, dbRules: make(config.DatabaseRules), securityFunctions: make(config.SecurityFunctions), dbPrepQueryRules: make(config.DatabasePreparedQueries), crud: crud, adminMan: adminMan, jwt: jwtUtils.New()}
}

// GetInternalAccessToken returns the token that can be used internally by Space Cloud
func (m *Module) GetInternalAccessToken(ctx context.Context) (string, error) {
	return m.CreateToken(ctx, map[string]interface{}{
		"id":     utils.InternalUserID,
		"nodeId": m.nodeID,
		"role":   "SpaceCloud",
	})
}

// GetSCAccessToken returns the token that can be used to verify Space Cloud
func (m *Module) GetSCAccessToken(ctx context.Context) (string, error) {
	return m.CreateToken(ctx, map[string]interface{}{
		"id":   m.nodeID,
		"role": "SpaceCloud",
	})
}

// CreateToken generates a new JWT Token with the token claims
func (m *Module) CreateToken(ctx context.Context, tokenClaims model.TokenClaims) (string, error) {
	return m.jwt.CreateToken(ctx, tokenClaims)
}

// IsTokenInternal checks if the provided token is internally generated
func (m *Module) IsTokenInternal(ctx context.Context, token string) error {
	claims, err := m.jwt.ParseToken(ctx, token)
	if err != nil {
		return err
	}

	if idTemp, p := claims["id"]; p {
		if id, ok := idTemp.(string); ok && id == utils.InternalUserID {
			return nil
		}
	}
	return errors.New("token has not been created internally")
}
