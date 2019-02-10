package auth

import (
	"errors"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// Module is responsible for authentication and authorsation
type Module struct {
	sync.RWMutex
	rules  config.Crud
	secret string
	crud   *crud.Module
}

// Init creates a new instance of the auth object
func Init(crud *crud.Module) *Module {
	return &Module{rules: make(config.Crud), crud: crud}
}

// SetConfig set the rules adn secret key required by the auth block
func (m *Module) SetConfig(secret string, rules config.Crud) {
	m.Lock()
	defer m.Unlock()

	m.rules = rules
	m.secret = secret
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
		}
	}
	return nil, ErrRuleNotFound
}

// CreateToken generates a new JWT Token
func (m *Module) CreateToken(obj map[string]interface{}) (string, error) {
	m.RLock()
	defer m.RUnlock()

	claims := jwt.MapClaims{}
	for k, v := range obj {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// IsAuthenticated checks if the caller is authentic
func (m *Module) IsAuthenticated(token, dbType, col string, query utils.OperationType) (map[string]interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getRule(dbType, col, query)
	if err != nil {
		return nil, err
	}

	if rule.Rule == "allow" {
		return map[string]interface{}{}, nil
	}

	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("Invalid signing method type")
		}

		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		obj := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			obj[key] = val
		}

		return obj, nil
	}

	return nil, errors.New("AUTH: JWT token could not be verified")
}

// IsAuthorized checks if the caller is authorized to make the request
func (m *Module) IsAuthorized(dbType, col string, query utils.OperationType, args map[string]interface{}) error {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getRule(dbType, col, query)
	if err != nil {
		return err
	}

	switch rule.Rule {
	case "allow", "authenticated":
		return nil

	case "deny":
		return ErrIncorrectMatch

	case "match":
		return match(rule, args)

	case "and":
		return matchAnd(rule, args)

	case "or":
		return matchOr(rule, args)

	case "query":
		return matchQuery(rule, m.crud, args)

	default:
		return ErrIncorrectMatch
	}
}
