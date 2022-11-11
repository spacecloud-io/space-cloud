package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/utils"
)

// SendUnauthenticatedResponse sends a unauthenticated response to the user
func SendUnauthenticatedResponse(w http.ResponseWriter, r *http.Request, logger *zap.Logger, err error) {
	zapFields := []zap.Field{zap.String("url", r.URL.Path), zap.String("host", r.Host)}
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}

	logger.Error("Unable to authenticate request", zapFields...)
	_ = utils.SendErrorResponse(w, http.StatusUnauthorized, errors.New("unable to authenticate request"))
}

// Sign signs a token using the primary secret
func (a *App) Sign(claims map[string]interface{}) (string, error) {
	secret := a.getPrimarySecret()

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
func (a *App) Verify(token string) (map[string]interface{}, error) {
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
		secret, p := a.getSecretFromKID(kid.(string))
		if !p {
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
		claims, ok := tokenObj.Claims.(jwt.MapClaims)
		if !ok || !tokenObj.Valid {
			return nil, errors.New("invalid claims provided")
		}

		if err := claims.Valid(); err != nil {
			return nil, err
		}

		obj := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			obj[key] = val
		}

		// Check if issuer is correct
		if len(secret.AllowedIssuers) > 0 {
			c, ok := claims["iss"]
			if !ok {
				return nil, errors.New("claim (iss) not provided in token")
			}
			if err := verifyClaims(c, secret.AllowedIssuers); err != nil {
				return nil, err
			}
		}

		// Check if audience is correct
		if len(secret.AllowedAudiences) > 0 {
			c, ok := claims["aud"]
			if !ok {
				return nil, errors.New("claim (aud) not provided in token")
			}
			if err := verifyClaims(c, secret.AllowedAudiences); err != nil {
				return nil, err
			}
		}

		return obj, nil
	}

	return nil, errors.New("unable to verify token")
}
