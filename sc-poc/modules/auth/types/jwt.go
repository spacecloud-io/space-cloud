package types

import "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"

type (
	// AuthSecret exposes a custom type made to manage authentication
	AuthSecret struct {
		v1alpha1.AuthSecret
		Alg   JWTAlg
		Value string
	}

	// JWTAlg is type of method used for signing token
	JWTAlg string
)

const (
	// HS256 is method used for signing token
	HS256 JWTAlg = "HS256"

	// RS256 is method used for signing token
	RS256 JWTAlg = "RS256"

	// JwkURL is the method for identifying a secret that has to be validated against secret kes fetched from url
	JwkURL JWTAlg = "JWK_URL"

	// RS256Public is the method for identifying a secret that has to be validated against with a public key
	RS256Public JWTAlg = "RS256_PUBLIC"
)
