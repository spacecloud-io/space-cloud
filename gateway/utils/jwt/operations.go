package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ParseToken verifies the token
func (j *JWT) ParseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	parser := jwt.Parser{}
	parsedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to parse token, token is malformed", err, nil)
	}
	kid, isKIDPresent := parsedToken.Header["kid"]
	if isKIDPresent {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Token kid", map[string]interface{}{"kid": kid})
		// check if kid belongs to a jwk secret
		claims, err := j.parseJwkSecret(ctx, kid.(string), token)
		if err == nil {
			return claims, nil
		}
		// check if kid belongs to a normal token with kid header
		obj, ok := j.staticSecrets[kid.(string)]
		if ok {
			if obj.Alg == config.RS256Public {
				obj.Alg = config.RS256
			}
			return j.verifyTokenSignature(ctx, token, obj)
		}
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "", nil, nil)
	}

	for _, secret := range j.staticSecrets {
		// normal token
		switch secret.Alg {
		case "":
			secret.Alg = config.HS256
		case config.RS256Public:
			secret.Alg = config.RS256
		}

		claims, err := j.verifyTokenSignature(ctx, token, secret)
		if err == nil {
			return claims, nil
		}
	}
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Authentication secrets not set in space cloud config", err, nil)
}

// CreateToken create a token with primary secret
func (j *JWT) CreateToken(ctx context.Context, tokenClaims model.TokenClaims) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}
	var tokenString string
	var err error
	// Add expiry of one week
	claims["exp"] = time.Now().Add(24 * 7 * time.Hour).Unix()
	for _, s := range j.staticSecrets {
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
				return "", helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid algorithm (%s) provided for creating token", s.Alg), err, nil)
			}
		}
	}
	return "", errors.New("no primary secret provided")
}
