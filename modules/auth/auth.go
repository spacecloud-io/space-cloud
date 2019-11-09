package auth

import (
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/schema"
	"github.com/spaceuptech/space-cloud/modules/crud"

	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/utils"
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
	secret          string
	crud            *crud.Module
	functions       *functions.Module
	fileRules       []*config.FileRule
	funcRules       *config.ServicesModule
	project         string
	fileStoreType   string
	Schema          *schema.Schema
	makeHttpRequest utils.MakeHttpRequest
}

// Init creates a new instance of the auth object
func Init(crud *crud.Module, functions *functions.Module, removeProjectScope bool) *Module {
	return &Module{rules: make(config.Crud), crud: crud, functions: functions, Schema: schema.Init(crud, removeProjectScope)}
}

// SetConfig set the rules and secret key required by the auth block
func (m *Module) SetConfig(project string, secret string, rules config.Crud, fileStore *config.FileStore, functions *config.ServicesModule) error {
	m.Lock()
	defer m.Unlock()

	if fileStore != nil {
		sortFileRule(fileStore.Rules)
	}

	m.project = project
	m.rules = rules

	if err := m.Schema.SetConfig(rules, project); err != nil {
		return err
	}
	m.secret = secret
	if fileStore != nil && fileStore.Enabled {
		m.fileRules = fileStore.Rules
		m.fileStoreType = fileStore.StoreType
	}

	if functions != nil {
		m.funcRules = functions
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
	return m.CreateToken(map[string]interface{}{"id": utils.InternalUserID})
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
