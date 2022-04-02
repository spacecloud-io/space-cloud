package project

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/managers/configman"
)

type (
	m map[string]interface{}
)

func getTypeDefinitions() configman.Types {
	return configman.Types{
		"config": &configman.TypeDefinition{
			Schema: m{
				"type": "object",
				"properties": m{
					"id": m{
						"type": "string",
					},
					"name": m{
						"type": "string",
					},
					"dockerRegistry": m{
						"type": "string",
					},
					"contextTimeGraphQL": m{
						"type": openapi3.TypeInteger,
					},
				},
			},
			Hooks:           configman.Hooks{},
			RequiredParents: []string{},
		},
		"aes-key": &configman.TypeDefinition{
			IsSecure: true,
			Schema: m{
				"type": "object",
				"properties": m{
					"key": m{
						"type": "string",
					},
				},
			},
			Hooks:           configman.Hooks{configman.PhasePreApply: struct{}{}},
			RequiredParents: []string{"project"},
		},
		"jwt-secret": &configman.TypeDefinition{
			IsSecure: true,
			Schema: m{
				"type": "object",
				"properties": m{
					"isPrimary": m{
						"type": openapi3.TypeBoolean,
					},
					"alg": m{
						"type": openapi3.TypeString,
						"enum": []string{string(config.HS256), string(config.RS256), string(config.JwkURL), string(config.RS256Public)},
					},
					"kid": m{
						"type": openapi3.TypeString,
					},
					"jwkUrl": m{
						"type": openapi3.TypeString,
					},
					"aud": m{
						"type":  openapi3.TypeArray,
						"items": m{"type": openapi3.TypeString},
					},
					"iss": m{
						"type":  openapi3.TypeArray,
						"items": m{"type": openapi3.TypeString},
					},
					"secret": m{
						"type": openapi3.TypeString,
					},
					"publicKey": m{
						"type": openapi3.TypeString,
					},
					"privateKey": m{
						"type": openapi3.TypeString,
					},
				},
			},
			Hooks:           configman.Hooks{},
			RequiredParents: []string{"project"},
		},
	}
}
