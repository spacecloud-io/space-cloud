package auth

import (
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for authentication and authorisation
type Module struct {
	sync.RWMutex
	rules     config.Crud
	secret    string
	crud      *crud.Module
	functions *functions.Module
	fileRules map[string]*config.FileRule
	funcRules config.FuncRules
	project   string
}

// TokenClaims holds the JWT token claims
type TokenClaims map[string]interface{}

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
	}

	if functions != nil && functions.Enabled {
		m.funcRules = functions.Rules
	}
}

// SetSecret sets the secret key to be used for JWT authentication
func (m *Module) SetSecret(secret string) {
	m.Lock()
	defer m.Unlock()
	m.secret = secret
}

func (m *Module) getRule(dbType, col string, query utils.OperationType) (*config.Rule, error) {
	if dbRules, p1 := m.rules[dbType]; p1 {
		if collection, p2 := dbRules.Collections[col]; p2 {
			if rule, p3 := collection.Rules[string(query)]; p3 {
				return rule, nil
			}
		} else if defaultCol, p2 := dbRules.Collections["default"]; p2 {
			if rule, p3 := defaultCol.Rules[string(query)]; p3 {
				return rule, nil
			}
		}
	}
	return nil, ErrRuleNotFound
}

// CreateToken generates a new JWT Token
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
		tokenClaims := make(TokenClaims, len(claims))
		for key, val := range claims {
			tokenClaims[key] = val
		}

		return tokenClaims, nil
	}

	return nil, errors.New("AUTH: JWT token could not be verified")
}

// IsAuthenticated checks if the caller is authentic
func (m *Module) IsAuthenticated(token, dbType, col string, query utils.OperationType) (TokenClaims, error) {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getRule(dbType, col, query)
	if err != nil {
		return nil, err
	}

	if rule.Rule == "allow" {
		return nil, nil
	}
	return m.parseToken(token)
}

// IsAuthorized checks if the caller is authorized to make the request
func (m *Module) IsAuthorized(project, dbType, col string, query utils.OperationType, args map[string]interface{}) error {
	m.RLock()
	defer m.RUnlock()

	if project != m.project {
		return errors.New("invalid project details provided")
	}

	rule, err := m.getRule(dbType, col, query)
	if err != nil {
		return err
	}

	return m.matchRule(project, rule, args)
}
