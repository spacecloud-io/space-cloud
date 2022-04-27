package graphql

import "github.com/graphql-go/graphql"

func getDirectors() []*graphql.Directive {
	return []*graphql.Directive{
		{
			Name:        "export",
			Description: "Export the value to use for on-the-fly joins",
			Args:        []*graphql.Argument{{PrivateName: "as", Type: graphql.String}},
			Locations:   []string{graphql.DirectiveLocationField},
		},
	}
}
