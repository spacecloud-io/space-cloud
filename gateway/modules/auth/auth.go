package auth

import (
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

var (
	// ErrInvalidSigningMethod denotes a jwt signing method type not used by Space Cloud.
	ErrInvalidSigningMethod = errors.New("invalid signing method type")
	ErrTokenVerification    = errors.New("AUTH: JWT token could not be verified")
)

// Module is responsible for authentication and authorisation
type Module struct {
	sync.RWMutex
	rules           config.Crud
	nodeID          string
	secret          string
	crud            *crud.Module
	fileRules       []*config.FileRule
	funcRules       *config.ServicesModule
	eventingRules   map[string]*config.Rule
	project         string
	fileStoreType   string
	schema          *schema.Schema
	makeHttpRequest utils.MakeHttpRequest
}

// PostProcess is responsible for implementing force and remove rules
type PostProcess struct {
	postProcessAction []PostProcessAction
}

// PostProcessAction has action ->  force/remove and field,value depending on the Action.
type PostProcessAction struct {
	Action string
	Field  string
	Value  interface{}
}

// Init creates a new instance of the auth object
func Init(nodeID string, crud *crud.Module, schema *schema.Schema, removeProjectScope bool) *Module {
	return &Module{nodeID: nodeID, rules: make(config.Crud), crud: crud, schema: schema}
}

// SetConfig set the rules and secret key required by the auth block
func (m *Module) SetConfig(project string, secret string, rules config.Crud, fileStore *config.FileStore, functions *config.ServicesModule, eventing *config.Eventing) error {
	m.Lock()
	defer m.Unlock()

	if fileStore != nil {
		sortFileRule(fileStore.Rules)
	}

	m.project = project
	m.rules = rules
	m.secret = secret
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

// SetSecret sets the secret key to be used for JWT authentication
func (m *Module) SetSecret(secret string) {
	m.Lock()
	defer m.Unlock()
	m.secret = secret
}

// GetInternalAccessToken returns the token that can be used internally by Space Cloud
func (m *Module) GetInternalAccessToken() (string, error) {
	return m.CreateToken(map[string]interface{}{
		"id":     utils.InternalUserID,
		"nodeId": m.nodeID,
		"role":   "SpaceCloud",
	})
}

// GetSCAccessToken returns the token that can be used to verify Space Cloud
func (m *Module) GetSCAccessToken() (string, error) {
	return m.CreateToken(map[string]interface{}{
		"id":   m.nodeID,
		"role": "SpaceCloud",
	})
}

// CreateToken generates a new JWT Token with the token claims
func (m *Module) CreateToken(tokenClaims TokenClaims) (string, error) {
	m.RLock()
	defer m.RUnlock()

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IsTokenInternal checks if the provided token is internally generated
func (m *Module) IsTokenInternal(token string) error {
	m.RLock()
	defer m.RUnlock()

	claims, err := m.parseToken(token)
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

func (m *Module) parseToken(token string) (TokenClaims, error) {
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		obj := make(TokenClaims, len(claims))
		for key, val := range claims {
			obj[key] = val
		}

		return obj, nil
	}

	return nil, ErrTokenVerification
}

func (m *Module) SetMakeHttpRequest(function utils.MakeHttpRequest) {
	m.Lock()
	defer m.Unlock()

	m.makeHttpRequest = function
}
