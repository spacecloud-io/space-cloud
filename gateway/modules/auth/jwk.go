package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

type jsonWebKeySet struct {
	Closer    chan struct{}
	TimeStamp int // time in seconds
	Set       *jwk.Set
}

func (m *Module) getJWKSet(id string) (*jsonWebKeySet, bool) {
	value, ok := m.jsonWebKeys[id]
	return value, ok
}

func (m *Module) fetchJWKKeys(ctx context.Context, url string) (*jwk.Set, http.Header, error) {
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

func (m *Module) getJWKRefreshTime(secret *config.Secret) (*jsonWebKeySet, error) {
	ctx := context.Background()

	obj := new(jsonWebKeySet)
	set, headers, err := m.fetchJWKKeys(ctx, secret.JwkURL)
	if err != nil {
		return nil, err
	}
	obj.Set = set
	for _, key := range set.Keys {
		fmt.Println("fuck", key.KeyID())
	}
	// check cache-control header for refresh time of jwks
	values := headers.Get("cache-control")
	if values != "" {
		var cacheTime string
		for _, value := range strings.Split(values, ",") {
			value = strings.TrimSpace(value)
			if strings.HasPrefix(value, "max-age") {
				cacheTime = strings.Split(value, "=")[1]
				break
			}
			if strings.HasPrefix(value, "s-maxage") {
				cacheTime = strings.Split(value, "=")[1]
				break
			}
		}
		obj.TimeStamp, err = strconv.Atoi(cacheTime)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwt url (%s), Cache-control header contains data of inavlid type expecting string", secret.JwkURL), nil, nil)
		}
		return obj, nil
	}

	// check expires header for refresh time of jwks
	values = headers.Get("expires")
	if values != "" {
		t, err := http.ParseTime(values)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to process jwt url (%s), Expires header contains data of inavlid type expecting string in format RFC1123", secret.JwkURL), nil, nil)
		}
		obj.TimeStamp = t.Second()
		return obj, nil
	}

	// set default refresh time
	obj.TimeStamp = time.Now().Add(24 * time.Hour).Second()
	return obj, nil
}
