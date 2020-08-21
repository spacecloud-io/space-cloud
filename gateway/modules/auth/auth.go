package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

var (
	// ErrInvalidSigningMethod denotes a jwt signing method type not used by Space Cloud.
	ErrInvalidSigningMethod = errors.New("invalid signing method type")
	// ErrTokenVerification is thrown when a jwt token could not be verified
	ErrTokenVerification = errors.New("AUTH: JWT token could not be verified")
)

// Module is responsible for authentication and authorisation
type Module struct {
	sync.RWMutex
	rules           config.Crud
	nodeID          string
	secrets         []*config.Secret
	crud            model.CrudAuthInterface
	fileRules       []*config.FileRule
	funcRules       *config.ServicesModule
	eventingRules   map[string]*config.Rule
	project         string
	fileStoreType   string
	makeHTTPRequest utils.TypeMakeHTTPRequest
	aesKey          []byte

	// Admin Manager
	adminMan       adminMan
	integrationMan integrationManagerInterface
}

// Init creates a new instance of the auth object
func Init(nodeID string, crud model.CrudAuthInterface, adminMan adminMan, integrationMan integrationManagerInterface) *Module {
	return &Module{nodeID: nodeID, rules: make(config.Crud), crud: crud, adminMan: adminMan, integrationMan: integrationMan}
}

// SetConfig set the rules and secret key required by the auth block
func (m *Module) SetConfig(project, secretSource string, secrets []*config.Secret, encodedAESKey string, rules config.Crud, fileStore *config.FileStore, functions *config.ServicesModule, eventing *config.Eventing) error {
	m.Lock()
	defer m.Unlock()

	if fileStore != nil {
		sortFileRule(fileStore.Rules)
	}

	m.project = project
	m.rules = rules

	// Store the secret
	m.secrets = secrets
	if secretSource == "admin" {
		m.secrets = []*config.Secret{{Secret: m.adminMan.GetSecret(), IsPrimary: true, Alg: config.HS256}}
	}

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

// SetSecrets sets the secrets to be used for JWT authentication
func (m *Module) SetSecrets(secretSource string, secrets []*config.Secret) {
	m.Lock()
	defer m.Unlock()
	m.secrets = secrets
	if secretSource == "admin" {
		m.secrets = []*config.Secret{{Secret: m.adminMan.GetSecret(), IsPrimary: true}}
	}
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

	sortFileRule(fileStore.Rules)
	m.fileRules = fileStore.Rules
	m.fileStoreType = fileStore.StoreType

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
func (m *Module) CreateToken(tokenClaims model.TokenClaims) (string, error) {
	m.RLock()
	defer m.RUnlock()

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}
	var tokenString string
	var err error
	// Add expiry of one week
	claims["exp"] = time.Now().Add(24 * 7 * time.Hour).Unix()
	for _, s := range m.secrets {
		if s.IsPrimary {
			switch s.Alg {
			case config.RS256:
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(s.PrivateKey))
				if err != nil {
					return "", err
				}
				tokenString, err = token.SignedString(signKey)
				if err != nil {
					return "", err
				}
				return tokenString, nil
			case config.HS256, "":
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenString, err = token.SignedString([]byte(s.Secret))
				if err != nil {
					return "", err
				}
				return tokenString, nil
			default:
				return "", fmt.Errorf("invalid algorithm (%s) provided", s.Alg)
			}
		}
	}
	return "", errors.New("no primary secret provided")
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

func (m *Module) parseToken(token string) (map[string]interface{}, error) {
	for _, secret := range m.secrets {
		if string(secret.Alg) == "" {
			secret.Alg = config.HS256
		}
		// Parse the JWT token
		tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if token.Method.Alg() != string(secret.Alg) {
				return nil, fmt.Errorf("invalid signing type provided - wanted %s; got %s", secret.Alg, token.Method.Alg())
			}
			switch secret.Alg {
			case config.RS256:
				return jwt.ParseRSAPublicKeyFromPEM([]byte(secret.PublicKey))
			case config.HS256, "":
				return []byte(secret.Secret), nil
			default:
				return nil, fmt.Errorf("invalid algorithm (%s) provided", secret.Alg)
			}
		})
		if err != nil {
			utils.LogDebug("Unable to verify jwt token", "admin", "parse-token", map[string]interface{}{"token": token, "error": err.Error()})
			continue
		}

		// Get the claims
		if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
			if err := claims.Valid(); err != nil {
				return nil, err
			}
			obj := make(TokenClaims, len(claims))
			for key, val := range claims {
				obj[key] = val
			}

			return obj, nil
		}
		utils.LogDebug("Unable to verify jwt token", "admin", "parse-token", map[string]interface{}{"token": token, "valid": tokenObj.Valid})
	}
	return nil, ErrTokenVerification
}

// SetMakeHTTPRequest sets the http request
func (m *Module) SetMakeHTTPRequest(function utils.TypeMakeHTTPRequest) {
	m.Lock()
	defer m.Unlock()

	m.makeHTTPRequest = function
}
