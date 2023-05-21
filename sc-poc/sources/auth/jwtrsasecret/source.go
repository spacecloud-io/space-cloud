package jwtrsasecret

import (
	"crypto/rsa"

	"github.com/caddyserver/caddy/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/auth"
	"github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var jwtrsasecretResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "jwtrsasecrets"}

func init() {
	source.RegisterSource(JWTRSASecretSource{}, jwtrsasecretResource)
}

// JWTRSASecretSource describes a JWTRSASecret source
type JWTRSASecretSource struct {
	v1alpha1.JwtRSASecret
	parsedPubKey *rsa.PublicKey
	parsedPriKey *rsa.PrivateKey

	// For internal use
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (JWTRSASecretSource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(jwtrsasecretResource)),
		New: func() caddy.Module { return new(JWTRSASecretSource) },
	}
}

// Provision provisions the source
func (s *JWTRSASecretSource) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s)

	var err error
	s.parsedPubKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(s.Spec.PublicKey.Value))
	if err != nil {
		s.logger.Error("Unable to parse rsa public key", zap.Error(err))
		return err
	}

	if s.Spec.PrivateKey != nil {
		s.parsedPriKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(s.Spec.PrivateKey.Value))
		if err != nil {
			s.logger.Error("Unable to parse rsa private key", zap.Error(err))
			return err
		}
	}

	return nil
}

// GetPriority returns the priority of the source.
func (s *JWTRSASecretSource) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *JWTRSASecretSource) GetProviders() []string {
	return []string{"auth"}
}

// GetSecretInfo returns the details of the secret such as
// algorithm used, kid, allowed audiences, allowed issuers
func (s *JWTRSASecretSource) GetSecretInfo() *types.AuthSecret {
	return &types.AuthSecret{
		AuthSecret: s.Spec.AuthSecret,
		Alg:        types.RS256,
		PublicKey:  s.parsedPubKey,
		PrivateKey: s.parsedPriKey,
	}
}

// Interface guard
var (
	_ source.Source     = (*JWTRSASecretSource)(nil)
	_ auth.SecretSource = (*JWTRSASecretSource)(nil)
)
