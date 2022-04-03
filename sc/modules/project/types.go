package project

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
)

type (
	m map[string]interface{}
)

func getTypeDefinitions() model.Types {
	return model.Types{
		"config": &model.TypeDefinition{
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
			Hooks:           model.Hooks{},
			RequiredParents: []string{},
		},
		"aes-key": &model.TypeDefinition{
			IsSecure: true,
			Schema: m{
				"type": "object",
				"properties": m{
					"key": m{
						"type": "string",
					},
				},
			},
			Hooks:           model.Hooks{model.PhasePreApply: struct{}{}},
			RequiredParents: []string{"project"},
		},
		"jwt-secret": &model.TypeDefinition{
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
			Hooks:           model.Hooks{},
			RequiredParents: []string{"project"},
		},
	}
}
