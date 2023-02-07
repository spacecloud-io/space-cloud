package auth

import (
	"context"

	"github.com/spacecloud-io/space-cloud/modules/auth/types"
)

type (
	// SecretSource describes the implementation of a secret source
	SecretSource interface {
		GetSecretInfo() *types.AuthSecret
	}

	// PolicySource describes the implementation of a policy source
	PolicySource interface {
		Evaluate(context.Context, interface{}) (bool, string, error)
	}
)
