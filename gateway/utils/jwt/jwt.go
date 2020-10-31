package jwt

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// JWT holds the required fields for jwt package
type JWT struct {
	lock                 sync.RWMutex
	staticSecrets        map[string]*config.Secret
	jwkSecrets           map[string]*jwkSecret
	closeJwkRoutineChan  chan struct{}
	mapJwkKidToSecretKid map[string]string
}

type jwkSecret struct {
	audience    []string
	issuer      []string
	url         string
	refreshTime time.Time
	set         *jwk.Set
}

const defaultRefreshTime = 1 * time.Hour

// New initializes the package
func New() *JWT {
	ch := make(chan struct{}, 1)
	j := &JWT{
		staticSecrets:        map[string]*config.Secret{},
		jwkSecrets:           map[string]*jwkSecret{},
		mapJwkKidToSecretKid: map[string]string{},
		closeJwkRoutineChan:  ch,
	}
	go func() {
		tick := time.NewTicker(defaultRefreshTime)
		for {
			select {
			case t := <-tick.C:
				j.fetchJWKRoutine(t)
			case <-ch:
				return
			}
		}
	}()
	return j
}

// SetSecrets set the internal fields of jwt struct
func (j *JWT) SetSecrets(secrets []*config.Secret) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	newJwkSecrets := map[string]*jwkSecret{}
	newKidMap := map[string]string{}
	newStaticSecretMap := map[string]*config.Secret{}
	for _, secret := range secrets {
		switch secret.Alg {
		case config.JwkURL:
			// Set the secret kid if it isn't already set
			if secret.KID == "" {
				secret.KID = ksuid.New().String()
			}

			// add or update a secret
			obj, ok := j.jwkSecrets[secret.KID]
			if !ok || secret.JwkURL != obj.url {
				jwkSecretInfo, err := getJWKRefreshTime(secret.JwkURL)
				if err != nil {
					return err
				}
				jwkSecretInfo.audience = secret.Audience
				jwkSecretInfo.issuer = secret.Issuer

				newJwkSecrets[secret.KID] = jwkSecretInfo
				for _, key := range jwkSecretInfo.set.Keys {
					newKidMap[key.KeyID()] = secret.KID
				}
				continue
			}

			// add existing secret to new map
			obj.audience = secret.Audience
			obj.issuer = secret.Issuer
			newJwkSecrets[secret.KID] = obj
			for key, value := range j.mapJwkKidToSecretKid {
				if value == secret.KID {
					newKidMap[key] = value
				}
			}

		case config.RS256Public:
			if secret.IsPrimary {
				return helpers.Logger.LogError("internal", "RSA algorithms without private keys cannot be used as a primary secret", nil, nil)
			}

			// Set the secret kid if it isn't already set
			if secret.KID == "" {
				h := sha256.New()
				_, _ = h.Write([]byte(secret.PublicKey))
				secret.KID = base64.StdEncoding.EncodeToString(h.Sum(nil))
			}

			newStaticSecretMap[secret.KID] = secret

		case config.RS256:
			// Set the secret kid if it isn't already set
			if secret.KID == "" {
				h := sha256.New()
				_, _ = h.Write([]byte(secret.PublicKey))
				secret.KID = base64.StdEncoding.EncodeToString(h.Sum(nil))
			}

			newStaticSecretMap[secret.KID] = secret
		case config.HS256, "":
			// Set the secret kid if it isn't already set
			if secret.KID == "" {
				h := sha256.New()
				_, _ = h.Write([]byte(secret.Secret))
				secret.KID = base64.StdEncoding.EncodeToString(h.Sum(nil))
			}

			newStaticSecretMap[secret.KID] = secret
		default:
			return helpers.Logger.LogError("internal", fmt.Sprintf("Invalid token algorithm provided (%s)", secret.Alg), nil, nil)
		}
	}

	j.jwkSecrets = newJwkSecrets
	j.staticSecrets = newStaticSecretMap
	j.mapJwkKidToSecretKid = newKidMap
	return nil
}
