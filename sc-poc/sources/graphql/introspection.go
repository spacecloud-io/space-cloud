package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"

	"github.com/spacecloud-io/space-cloud/modules/graphql/types"
)

func (s *GraphqlSource) getRawGraphqlSchema() error {
	data, _ := json.Marshal(map[string]string{"query": getIntrospectionQuery()})
	resp, err := http.Post(s.Spec.Source.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Extract the body
	rawSchema := new(introspectionResponse)
	_ = json.NewDecoder(resp.Body).Decode(rawSchema)
	s.rawSchema = rawSchema
	return nil
}

func (s *GraphqlSource) prepareGraphqlTypes(queryRootType, mutationRootType graphql.Fields, graphqlTypes map[string]graphql.Type) {
	s.getTypes(s.rawSchema.Data.Schema.Types, graphqlTypes)

	// TODO: Also accomodate the directives

	s.addToRootType(s.rawSchema.Data.Schema.QueryType.Name, queryRootType, graphqlTypes)
	s.addToRootType(s.rawSchema.Data.Schema.MutationType.Name, mutationRootType, graphqlTypes)
}

func (s *GraphqlSource) getTypes(types []*introspectionResponseType, graphqlTypes map[string]graphql.Type) {
	// First we populate all the types
	for _, t := range types {
		name := getGraphqlTypeName(s.Name, t)
		if _, p := graphqlTypes[name]; p {
			continue
		}

		if v := s.getGraphqlType(s.Name, t); v != nil {
			graphqlTypes[name] = v
		}
	}

	// We now have to populate the fields for types `OBJECT` & `INPUT_OBJECT`
	for _, t := range types {
		s.populateGraphqlFields(t, graphqlTypes)
	}
}

func (s *GraphqlSource) getGraphqlType(srcName string, t *introspectionResponseType) graphql.Type {
	switch t.Kind {
	case graphql.TypeKindScalar:
		switch t.Name {
		case graphql.ID.Name(), graphql.Float.Name(), graphql.Int.Name(),
			graphql.String.Name(), graphql.Boolean.Name(), graphql.DateTime.Name():
			return nil
		}

		return graphql.NewScalar(graphql.ScalarConfig{
			Name:        t.Name,
			Description: t.Description,
		})

	case graphql.TypeKindEnum:
		enumMap := make(graphql.EnumValueConfigMap, len(t.EnumValues))
		for _, enumValue := range t.EnumValues {
			enumMap[enumValue.Name] = &graphql.EnumValueConfig{
				Description:       enumValue.Description,
				DeprecationReason: enumValue.DeprecationReason,
			}
		}
		return graphql.NewEnum(graphql.EnumConfig{
			Name:        getGraphqlTypeName(srcName, t),
			Description: t.Description,
			Values:      enumMap,
		})

	case graphql.TypeKindInputObject:
		// We will populate the fields later
		fields := make(graphql.InputObjectConfigFieldMap, len(t.InputFields))
		return graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        getGraphqlTypeName(srcName, t),
			Description: t.Description,
			Fields:      fields,
		})

	case graphql.TypeKindObject:
		// We will populate the fields later
		fields := make(graphql.Fields, len(t.Fields))
		return graphql.NewObject(graphql.ObjectConfig{
			Name:        getGraphqlTypeName(srcName, t),
			Description: t.Description,
			Fields:      fields,
		})

	case graphql.TypeKindInterface:
		// We don't support this at the moment
		return nil
	}

	return nil
}

func (s *GraphqlSource) populateGraphqlFields(t *introspectionResponseType, graphqlTypes map[string]graphql.Type) {
	switch t.Kind {
	case graphql.TypeKindInputObject:
		// Get the stored type
		storedType := graphqlTypes[getGraphqlTypeName(s.Name, t)].(*graphql.InputObject)

		// Add the fields to our type
		for _, field := range t.InputFields {
			storedType.AddFieldConfig(field.Name, &graphql.InputObjectFieldConfig{
				Type:         s.evaluateTypeRef(s.Name, field.TypeRef, graphqlTypes),
				Description:  field.Description,
				DefaultValue: field.DefaultValue,
			})
		}
	case graphql.TypeKindObject:
		// Get the stored type
		storedType := graphqlTypes[getGraphqlTypeName(s.Name, t)].(*graphql.Object)
		// Add the fields to our type
		for _, field := range t.Fields {

			// Prepare the arguments
			args := make(graphql.FieldConfigArgument, len(field.Args))
			for _, arg := range field.Args {
				args[arg.Name] = &graphql.ArgumentConfig{
					Type:         s.evaluateTypeRef(s.Name, arg.TypeRef, graphqlTypes),
					Description:  arg.Description,
					DefaultValue: arg.DefaultValue,
				}
			}

			// Add the field
			storedType.AddFieldConfig(field.Name, &graphql.Field{
				Type:              s.evaluateTypeRef(s.Name, field.TypeRef, graphqlTypes),
				Args:              args,
				Description:       field.Description,
				DeprecationReason: field.DeprecationReason,
			})
		}
	}
}

func (s *GraphqlSource) evaluateTypeRef(srcName string, gt *introspectionResponseTypeRef, graphqlTypes map[string]graphql.Type) graphql.Type {
	switch gt.Kind {
	case graphql.TypeKindList:
		t := s.evaluateTypeRef(srcName, gt.OfType, graphqlTypes)
		return graphql.NewList(t)

	case graphql.TypeKindNonNull:
		t := s.evaluateTypeRef(srcName, gt.OfType, graphqlTypes)
		return graphql.NewNonNull(t)

	default:
		return graphqlTypes[getGraphqlTypeName(srcName, gt)]
	}
}

func getGraphqlTypeName(srcName string, gt types.GraphqlType) string {
	// Prepend src name to prevent conflicts between multiple sources
	name := fmt.Sprintf("%s_%s", srcName, gt.GetName())

	// We don't want to modify names for scalar types
	if gt.GetKind() == graphql.TypeKindScalar {
		name = gt.GetName()
	}

	// Return the final name
	return name
}

func getIntrospectionQuery() string {
	return `
  query IntrospectionQuery {
    __schema {
      queryType { name }
      mutationType { name }
      subscriptionType { name }
      types {
        ...FullType
      }
      directives {
        name
        description
        locations
        args {
          ...InputValue
        }
      }
    }
  }
  fragment FullType on __Type {
    kind
    name
    description
    fields(includeDeprecated: true) {
      name
      description
      args {
        ...InputValue
      }
      type {
        ...TypeRef
      }
      isDeprecated
      deprecationReason
    }
    inputFields {
      ...InputValue
    }
    interfaces {
      ...TypeRef
    }
    enumValues(includeDeprecated: true) {
      name
      description
      isDeprecated
      deprecationReason
    }
    possibleTypes {
      ...TypeRef
    }
  }
  fragment InputValue on __InputValue {
    name
    description
    type { ...TypeRef }
    defaultValue
  }
  fragment TypeRef on __Type {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                }
              }
            }
          }
        }
      }
    }
  }
  `
}
