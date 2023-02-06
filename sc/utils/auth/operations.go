package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// Close shuts down all module operations
func (m *Module) Close() error {
	return nil
}

// Sign signs a token using the primary secret
func (m *Module) Sign(claims map[string]interface{}) (string, error) {
	secret := m.getPrimarySecret()

	// Prepare jwt claims
	jwtClaims := make(jwt.MapClaims, len(claims))
	for k, v := range claims {
		jwtClaims[k] = v
	}

	// Sign the token
	token := jwt.NewWithClaims(getSigninMethod(secret.Alg), jwtClaims)
	if secret.KID != "" {
		token.Header["kid"] = secret.KID
	}

	return token.SignedString(getSigningSecret(secret))
}

// Verify checks the validity of the token
func (m *Module) Verify(token string) (map[string]interface{}, error) {
	// Parse the token first
	parser := jwt.NewParser()
	parsedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	// Find the right secret to use
	kid, isKIDPresent := parsedToken.Header["kid"]
	if isKIDPresent {
		// Find the secret based on the kid
		secret, p := m.getSecretFromKID(kid.(string))
		if p {
			return nil, fmt.Errorf("no secret found using kid '%s'", kid)
		}

		tokenObj, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if t.Method.Alg() != string(secret.Alg) {
				return nil, fmt.Errorf("invalid token algorithm provided wanted (%s) got (%s)", secret.Alg, t.Method.Alg())
			}

			return getVerifyingSecret(secret), nil
		})
		if err != nil {
			return nil, err
		}

		// Get the claims
		if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
			if err := claims.Valid(); err != nil {
				return nil, err
			}
			obj := make(map[string]interface{}, len(claims))
			for key, val := range claims {
				obj[key] = val
			}

			// Check if issuer is correct
			if len(secret.Issuer) > 0 {
				c, ok := claims["iss"]
				if !ok {
					return nil, errors.New("claim (iss) not provided in token")
				}
				if err := verifyClaims(c, secret.Issuer); err != nil {
					return nil, err
				}
			}

			// Check if audience is correct
			if len(secret.Audience) > 0 {
				c, ok := claims["aud"]
				if !ok {
					return nil, errors.New("claim (aud) not provided in token")
				}
				if err := verifyClaims(c, secret.Audience); err != nil {
					return nil, err
				}
			}
		}
	}

	return nil, errors.New("unable to verify token")
}
