package database

import "github.com/spacecloud-io/space-cloud/model"

// GetConfigTypes returns all the config types expsed by the database app
func (a *App) GetConfigTypes() model.ConfigTypes {
	return a.getTypeDefinitions()
}

func (a *App) getTypeDefinitions() model.ConfigTypes {
	return model.ConfigTypes{
		"config": &model.ConfigTypeDefinition{
			Schema: m{
				"type": "object",
				"properties": m{
					"dbAlias": m{
						"type": "string",
					},
					"type": m{
						"type": "string",
					},
					"name": m{
						"type": "string",
					},
					"conn": m{
						"type": "string",
					},
					"isPrimary": m{
						"type": "boolean",
					},
					"enabled": m{
						"type": "boolean",
					},
					"batchTime": m{
						"type": "integer",
					},
					"batchRecords": m{
						"type": "integer",
					},
					"limit": m{
						"type": "integer",
					},
					"driverConf": m{
						"type": "object",
						"properties": m{
							"maxConn": m{
								"type": "integer",
							},
							"maxIdleTimeout": m{
								"type": "integer",
							},
							"minConn": m{
								"type": "integer",
							},
							"maxIdleConn": m{
								"type": "integer",
							},
						},
						"required": t{"maxConn", "maxIdleTimeout", "minConn", "maxIdleConn"},
					},
				},
				"required": t{"type", "name", "conn"},
			},
			RequiredParents: []string{"project"},
			Controller:      model.ConfigHooks{PreApply: processConfigHook},
		},
		"schema": &model.ConfigTypeDefinition{
			Schema: m{
				"type": "object",
				"properties": m{
					"col": m{
						"type": "string",
					},
					"dbAlias": m{
						"type": "string",
					},
					"schema": m{
						"type": "string",
					},
				},
				"required": t{"schema"},
			},
			RequiredParents: []string{"project", "database"},
			Controller:      model.ConfigHooks{PreApply: a.processDBSchemaHook},
		},
		"prepared-query": &model.ConfigTypeDefinition{
			Schema: m{
				"type": "object",
				"properties": m{
					"id": m{
						"type": "string",
					},
					"sql": m{
						"type": "string",
					},
					"rule": m{
						"type":                 "object",
						"additionalProperties": true,
					},
					"dbAlias": m{
						"type": "string",
					},
					"args": m{
						"type": "array",
						"items": m{
							"type": "string",
						},
					},
				},
				"required": t{"sql"},
			},
			RequiredParents: []string{"project", "database"},
			Controller: model.ConfigHooks{
				PreApply: processPreparedQuery,
			},
		},
	}
}
