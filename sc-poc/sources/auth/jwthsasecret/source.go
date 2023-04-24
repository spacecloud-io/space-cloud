package jwthsasecret

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/auth"
	"github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var jwthsasecretResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "jwthsasecrets"}

func init() {
	source.RegisterSource(JWTHSASecretSource{}, jwthsasecretResource)
}

// JWTHSASecretSource describes a JWTHSASecret source
type JWTHSASecretSource struct {
	v1alpha1.JwtHSASecret
}

// CaddyModule returns the Caddy module information.
func (JWTHSASecretSource) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(jwthsasecretResource)),
		New: func() caddy.Module { return new(JWTHSASecretSource) },
	}
}

// GetPriority returns the priority of the source.
func (s *JWTHSASecretSource) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *JWTHSASecretSource) GetProviders() []string {
	return []string{"auth"}
}

// GetSecretInfo returns the details of the secret such as
// algorithm used, kid, allowed audiences, allowed issuers
func (s *JWTHSASecretSource) GetSecretInfo() *types.AuthSecret {
	return &types.AuthSecret{
		AuthSecret: s.Spec.AuthSecret,
		Alg:        types.HS256,
		PublicKey:  s.Spec.Value,
		PrivateKey: s.Spec.Value,
	}
}

// Interface guard
var (
	_ source.Source     = (*JWTHSASecretSource)(nil)
	_ auth.SecretSource = (*JWTHSASecretSource)(nil)
)
