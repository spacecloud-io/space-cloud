package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ParseToken verifies the token
func (j *JWT) ParseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

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
			tempSecret := *obj
			switch obj.Alg {
			case "":
				tempSecret.Alg = config.HS256
			case config.RS256Public:
				tempSecret.Alg = config.RS256
			}
			return j.verifyTokenSignature(ctx, token, &tempSecret)
		}
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to parse token or secret with given kid does not exists", err, map[string]interface{}{"kid": kid})
	}

	var er error
	for _, secret := range j.staticSecrets {
		tempSecret := *secret
		// normal token
		switch secret.Alg {
		case "":
			tempSecret.Alg = config.HS256
		case config.RS256Public:
			tempSecret.Alg = config.RS256
		}

		claims, err := j.verifyTokenSignature(ctx, token, &tempSecret)
		if err == nil {
			return claims, nil
		}
		er = err
	}
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to parse token or authentication secrets not set in space cloud config", er, nil)
}

// CreateToken create a token with primary secret
func (j *JWT) CreateToken(ctx context.Context, tokenClaims model.TokenClaims) (string, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}
	var tokenString string
	var err error
	// Add expiry of one week
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()
	for _, s := range j.staticSecrets {
		if s.IsPrimary {
			switch s.Alg {
			case config.RS256:
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				token.Header["kid"] = s.KID
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
				token.Header["kid"] = s.KID
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

// Close closes the go routine of jwk fetch routing
func (j *JWT) Close() {
	j.lock.Lock()
	defer j.lock.Unlock()

	j.closeJwkRoutineChan <- struct{}{}

	j.staticSecrets = map[string]*config.Secret{}
	j.jwkSecrets = map[string]*jwkSecret{}
	j.mapJwkKidToSecretKid = map[string]string{}
}
