package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/redis/go-redis/v9"

	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var redisPool = caddy.NewUsagePool()

type client struct {
	*redis.Client
}

func createKey(spec v1alpha1.RedisSourceSpec) string {
	data, _ := json.Marshal(spec)
	return string(data)
}

func createNewClient(ctx context.Context, spec v1alpha1.RedisSourceSpec) caddy.Constructor {
	return func() (caddy.Destructor, error) {
		var password string
		if spec.Password != nil {
			password = spec.Password.Value
		}

		c := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", spec.Host.Value, spec.Port.Value),
			Password: password,
			DB:       0, // use default DB
		})

		// Check if a connection was established
		if err := c.Ping(ctx).Err(); err != nil {
			return nil, err
		}

		return &client{c}, nil
	}
}

// Destruct is called by the caddy.UsagePool during wrap up
func (c *client) Destruct() error {
	return c.Client.Close()
}

// Interface guards
var (
	_ caddy.Destructor = (*client)(nil)
)
