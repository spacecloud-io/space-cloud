package auth

import (
	"errors"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/crud"
)

// Rule is the authorisation object received from space cloud
type Rule struct {
	Rule      string                 `json:"rule"`
	Eval      string                 `json:"eval"`
	FieldType string                 `json:"type"`
	F1        interface{}            `json:"f1"`
	F2        interface{}            `json:"f2"`
	Clauses   []*Rule                `json:"clauses"`
	DbType    string                 `json:"db"`
	Col       string                 `json:"col"`
	Find      map[string]interface{} `json:"find"`
}

// Module is responsible for authentication and authorsation
type Module struct {
	sync.RWMutex
	rules  map[string]map[string]*Rule
	secret string
	crud   *crud.Module
}

// Init creates a new instance of the auth object
func Init(crud *crud.Module) *Module {
	return &Module{rules: make(map[string]map[string]*Rule), crud: crud}
}

// SetSecret sets the secret key to be used for JWT authentication
func (m *Module) SetSecret(secret string) {
	m.Lock()
	defer m.Unlock()
	m.secret = secret
}

// GetRule returns the rule associated with the particular database and collection.
// It returns an ErrRuleNotFound when the rules doesn't exist
func (m *Module) GetRule(dbType, col string) (*Rule, error) {
	m.RLock()
	defer m.RUnlock()
	if dbRules, p1 := m.rules[dbType]; p1 {
		if rule, p2 := dbRules[col]; p2 {
			return rule, nil
		}
	}
	return nil, ErrRuleNotFound
}

// IsAuthenticated checks if the caller is authentic
func (m *Module) IsAuthenticated(token, dbType, col string) (map[string]interface{}, error) {
	rule, err := m.GetRule(dbType, col)
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
func (m *Module) IsAuthorized(dbType, col string, args map[string]interface{}) error {
	rule, err := m.GetRule(dbType, col)
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
		// TODO: case for query

	default:
		return ErrIncorrectMatch
	}
}
