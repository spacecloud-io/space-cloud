package auth

import (
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
)

// Module is responsible for authentication and authorisation
type Module struct {
	sync.RWMutex
	rules         config.Crud
	secret        string
	crud          *crud.Module
	functions     *functions.Module
	fileRules     []*config.FileRule
	funcRules     config.Services
	project       string
	fileStoreType string
}

// Init creates a new instance of the auth object
func Init(crud *crud.Module, functions *functions.Module) *Module {
	return &Module{rules: make(config.Crud), crud: crud, functions: functions}
}

// SetConfig set the rules and secret key required by the auth block
func (m *Module) SetConfig(project string, secret string, rules config.Crud, fileStore *config.FileStore, functions *config.Functions) {
	m.Lock()
	defer m.Unlock()

	m.project = project
	m.rules = rules
	m.secret = secret
	if fileStore != nil && fileStore.Enabled {
		m.fileRules = fileStore.Rules
		m.fileStoreType = fileStore.StoreType
	}

	if functions != nil && functions.Enabled {
		m.funcRules = functions.Services
	}
}

// SetSecret sets the secret key to be used for JWT authentication
func (m *Module) SetSecret(secret string) {
	m.Lock()
	defer m.Unlock()
	m.secret = secret
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

func (m *Module) parseToken(token string) (TokenClaims, error) {
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method type")
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

	return nil, errors.New("AUTH: JWT token could not be verified")
}
