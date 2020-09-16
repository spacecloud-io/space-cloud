package jwt

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// JWT holds the required fields for jwt package
type JWT struct {
	staticSecrets        map[string]*config.Secret
	jwkSecrets           map[string]*jwkSecret
	mapJwkKidToSecretKid map[string]string
}

type jwkSecret struct {
	url         string
	refreshTime time.Time
	set         *jwk.Set
}

const defaultRefreshTime = 1 * time.Hour

// New initializes the package
func New() *JWT {
	j := new(JWT)
	go func() {
		tick := time.NewTicker(defaultRefreshTime)
		for t := range tick.C {
			for kid, secret := range j.jwkSecrets {
				if secret.refreshTime.After(t) {
					jwkSecretInfo, err := getJWKRefreshTime(secret.url)
					if err != nil {
						_ = helpers.Logger.LogError("", fmt.Sprintf("Unable to refresh jwk keys having kid (%s) and url (%s)", kid, secret.url), err, nil)
					}
					j.jwkSecrets[kid] = jwkSecretInfo
				}
			}
		}
	}()
	return j
}

// SetSecrets set the internal fields of jwt struct
func (j *JWT) SetSecrets(secrets []*config.Secret) error {
	newJwkSecrets := map[string]*jwkSecret{}
	newKidMap := map[string]string{}
	newStaticSecretMap := map[string]*config.Secret{}
	for _, secret := range secrets {
		if secret.KID == "" {
			secret.KID = ksuid.New().String()
		}
		switch secret.Alg {
		case config.JwkURL:
			// add or update a secret
			obj, ok := j.jwkSecrets[secret.KID]
			if !ok || secret.JwkURL != obj.url {
				jwkSecretInfo, err := getJWKRefreshTime(secret.JwkURL)
				if err != nil {
					return err
				}

				newJwkSecrets[secret.KID] = jwkSecretInfo
				for _, key := range jwkSecretInfo.set.Keys {
					newKidMap[key.KeyID()] = secret.KID
				}
				continue
			}

			// add existing secret to new map
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
			newStaticSecretMap[secret.KID] = secret

		case config.RS256, config.HS256, "":
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
