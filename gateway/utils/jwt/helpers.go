package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func fetchJWKKeys(ctx context.Context, url string) (*jwk.Set, http.Header, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to fetch jwks from provided url (%s)", url), err, nil)
	}
	if res.StatusCode != http.StatusOK {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwk url (%s), auth server returned status code (%v)", url, res.StatusCode), nil, nil)
	}
	set, err := jwk.Parse(res.Body)
	if err != nil {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwk url (%s), auth server has invalid response body", url), err, nil)
	}
	if set.Len() == 0 {
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwk url (%s), auth server returned 0 jwk keys", url), nil, nil)
	}
	return set, res.Header, nil
}

func getJWKRefreshTime(url string) (*jwkSecret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	obj := new(jwkSecret)
	set, headers, err := fetchJWKKeys(ctx, url)
	if err != nil {
		return nil, err
	}
	obj.set = set
	obj.url = url

	// check cache-control header for refresh time of jwks
	values := headers.Get("cache-control")
	if values != "" {
		var cacheTime string
		for _, value := range strings.Split(values, ",") {
			value = strings.TrimSpace(value)

			// Make default cache time 5 minutes
			if value == "no-cache" {
				cacheTime = strconv.Itoa(5 * 60)
				break
			}

			if strings.HasPrefix(value, "max-age") {
				cacheTime = strings.Split(value, "=")[1]
				break
			}
			if strings.HasPrefix(value, "s-maxage") {
				cacheTime = strings.Split(value, "=")[1]
				break
			}
		}
		duration, err := strconv.Atoi(cacheTime)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwt url (%s), Cache-control header contains data of inavlid type expecting string", url), err, nil)
		}
		obj.refreshTime = time.Now().Add(time.Duration(duration) * time.Second)
		return obj, nil
	}

	// check expires header for refresh time of jwks
	values = headers.Get("expires")
	if values != "" {
		t, err := http.ParseTime(values)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwt url (%s), Expires header contains data of inavlid type expecting string in format RFC1123", url), nil, nil)
		}
		obj.refreshTime = t
		return obj, nil
	}

	// set default refresh time
	obj.refreshTime = time.Now().Add(24 * time.Hour)
	return obj, nil
}

func (j *JWT) verifyTokenSignature(ctx context.Context, token string, secret *config.Secret) (map[string]interface{}, error) {
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != string(secret.Alg) {
			return nil, fmt.Errorf("invalid token algorithm provided wanted (%s) got (%s)", secret.Alg, token.Method.Alg())
		}

		if secret.JwkURL != "" {
			return secret.JwkKey, nil
		}

		switch secret.Alg {
		case config.RS256:
			return jwt.ParseRSAPublicKeyFromPEM([]byte(secret.PublicKey))
		case config.HS256, "":
			return []byte(secret.Secret), nil
		default:
			return nil, fmt.Errorf("invalid token algorithm (%s) provided", secret.Alg)
		}
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

		if len(secret.Issuer) > 0 {
			c, ok := claims["iss"]
			if !ok {
				return nil, errors.New("claim (iss) not provided in token")
			}
			if err := verifyClaims(c, secret.Issuer); err != nil {
				return nil, err
			}
		}

		if len(secret.Audience) > 0 {
			c, ok := claims["aud"]
			if !ok {
				return nil, errors.New("claim (aud) not provided in token")
			}
			if err := verifyClaims(c, secret.Audience); err != nil {
				return nil, err
			}
		}

		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Claim from request token", map[string]interface{}{"claims": claims, "type": "auth"})
		return obj, nil
	}
	return nil, errors.New("AUTH: JWT token could not be verified")
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

func (j *JWT) parseJwkSecret(ctx context.Context, kid, token string) (map[string]interface{}, error) {
	jwkSecretKID, ok := j.mapJwkKidToSecretKid[kid]
	if ok {
		obj, ok := j.jwkSecrets[jwkSecretKID]
		if !ok {
			return nil, errors.New("token contains an unknown kid")
		}
		key := obj.set.LookupKeyID(kid)
		if len(key) == 0 {
			return nil, fmt.Errorf("token has a invalid kid which doesn't match with the kid from the jwk url (%s)", obj.url)
		}
		var raw interface{}
		if err := key[0].Raw(&raw); err != nil {
			return nil, err
		}
		return j.verifyTokenSignature(ctx, token, &config.Secret{Issuer: obj.issuer, Audience: obj.audience, Alg: config.JWTAlg(key[0].Algorithm()), JwkKey: raw, JwkURL: obj.url})
	}
	return nil, errors.New("kid doesn't exists in internal jwk mapping of keys")
}

func (j *JWT) fetchJWKRoutine(t time.Time) {
	j.lock.Lock()
	defer j.lock.Unlock()

	for kid, secret := range j.jwkSecrets {
		if secret.refreshTime.Before(t) {
			jwkSecretInfo, err := getJWKRefreshTime(secret.url)
			if err != nil {
				_ = helpers.Logger.LogError("", fmt.Sprintf("Unable to refresh jwk keys having kid (%s) and url (%s)", kid, secret.url), err, nil)
				continue
			}
			jwkSecretInfo.audience = secret.audience
			jwkSecretInfo.issuer = secret.issuer
			j.jwkSecrets[kid] = jwkSecretInfo
			// Delete kids of existing JWK url
			for key, value := range j.mapJwkKidToSecretKid {
				if value == kid {
					delete(j.mapJwkKidToSecretKid, key)
				}
			}
			// Add new kids to existing JWK url
			for _, key := range jwkSecretInfo.set.Keys {
				j.mapJwkKidToSecretKid[key.KeyID()] = kid
			}
		}
	}
}
