package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"

	"github.com/spacecloud-io/space-cloud/config"
)

func getSigninMethod(alg config.JWTAlg) jwt.SigningMethod {
	switch alg {
	case config.HS256:
		return jwt.SigningMethodHS256
	case config.RS256, config.RS256Public, config.JwkURL:
		return jwt.SigningMethodRS256
	}

	return jwt.SigningMethodHS256
}

func getSigningSecret(secret *config.Secret) interface{} {
	if secret.Alg == config.HS256 {
		return []byte(secret.Secret)
	}

	// TODO: Return private key here
	return nil
}

func getVerifyingSecret(secret *config.Secret) interface{} {
	if secret.Alg == config.HS256 {
		return []byte(secret.Secret)
	}

	// TODO: Return public key here
	return nil
}

func (m *Module) getPrimarySecret() *config.Secret {
	for _, s := range m.secrets {
		if s.IsPrimary {
			return s
		}
	}

	return m.secrets[0]
}

func (m *Module) getSecretFromKID(kid string) (*config.Secret, bool) {
	for _, s := range m.secrets {
		if s.KID == kid {
			return s, true
		}
	}

	return nil, false
}

func verifyClaims(claims interface{}, configValues []string) error {
	switch claimFromToken := claims.(type) {
	case string:
		for _, cmp := range configValues {
			if claimFromToken == cmp {
				return nil
			}
		}
	case []interface{}:
		for _, audClaim := range claimFromToken {
			for _, c := range configValues {
				if audClaim == c {
					return nil
				}
			}
		}
	default:
		return errors.New("invalid type provided for claim in jwt token")
	}
	return errors.New("unable to verify claim of jwt token")
}
