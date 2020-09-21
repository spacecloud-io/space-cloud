package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/spaceuptech/space-cloud/gateway/utils"
	jwtUtils "github.com/spaceuptech/space-cloud/gateway/utils/jwt"
)

var (
	// ErrInvalidSigningMethod denotes a jwt signing method type not used by Space Cloud.
	ErrInvalidSigningMethod = errors.New("invalid signing method type")
)

// Module is responsible for authentication and authorisation
type Module struct {
	sync.RWMutex
	rules           config.Crud
	nodeID          string
	jwt             *jwtUtils.JWT
	crud            model.CrudAuthInterface
	fileRules       []*config.FileRule
	funcRules       *config.ServicesModule
	eventingRules   map[string]*config.Rule
	project         string
	fileStoreType   string
	makeHTTPRequest utils.TypeMakeHTTPRequest
	aesKey          []byte

	// Admin Manager
	adminMan adminMan
}

// Init creates a new instance of the auth object
func Init(nodeID string, crud model.CrudAuthInterface, adminMan adminMan) *Module {
	return &Module{nodeID: nodeID, rules: make(config.Crud), crud: crud, adminMan: adminMan, jwt: jwtUtils.New()}
}

// SetConfig set the rules and secret key required by the auth block
func (m *Module) SetConfig(project, secretSource string, secrets []*config.Secret, encodedAESKey string, rules config.Crud, fileStore *config.FileStore, functions *config.ServicesModule, eventing *config.Eventing) error {
	m.Lock()
	defer m.Unlock()

	if fileStore != nil {
		sortFileRule(fileStore.Rules)
	}

	if secretSource == "admin" {
		secrets = []*config.Secret{{KID: utils.AdminSecretKID, Secret: m.adminMan.GetSecret(), IsPrimary: true, Alg: config.HS256}}
	}
	if err := m.jwt.SetSecrets(secrets); err != nil {
		return err
	}

	m.project = project
	m.rules = rules

	decodedAESKey, err := base64.StdEncoding.DecodeString(encodedAESKey)
	if err != nil {
		return err
	}
	m.aesKey = decodedAESKey
	if fileStore != nil && fileStore.Enabled {
		m.fileRules = fileStore.Rules
		m.fileStoreType = fileStore.StoreType
	}

	if functions != nil {
		m.funcRules = functions
	}

	if eventing.SecurityRules != nil {
		m.eventingRules = eventing.SecurityRules
	}

	return nil
}

// SetAESKey sets the aeskey to be used for encryption
func (m *Module) SetAESKey(encodedAESKey string) error {
	m.Lock()
	defer m.Unlock()
	decodedAESKey, err := base64.StdEncoding.DecodeString(encodedAESKey)
	if err != nil {
		return err
	}
	m.aesKey = decodedAESKey
	return nil
}

// SetServicesConfig sets the service module config
func (m *Module) SetServicesConfig(projectID string, services *config.ServicesModule) {
	m.Lock()
	defer m.Unlock()
	m.project = projectID
	m.funcRules = services
}

// SetFileStoreConfig sets the file store module config
func (m *Module) SetFileStoreConfig(projectID string, fileStore *config.FileStore) {
	m.Lock()
	defer m.Unlock()
	m.project = projectID
	m.fileRules = fileStore.Rules
}

// SetEventingConfig sets the eventing config
func (m *Module) SetEventingConfig(securiyRules map[string]*config.Rule) {
	m.Lock()
	defer m.Unlock()
	m.eventingRules = securiyRules
}

// SetCrudConfig sets the crud module config
func (m *Module) SetCrudConfig(projectID string, crud config.Crud) {
	m.Lock()
	defer m.Unlock()
	m.project = projectID
	m.rules = crud
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

// SetMakeHTTPRequest sets the http request
func (m *Module) SetMakeHTTPRequest(function utils.TypeMakeHTTPRequest) {
	m.Lock()
	defer m.Unlock()

	m.makeHTTPRequest = function
}
