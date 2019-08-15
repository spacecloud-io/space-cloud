package graphql

import "testing"

func TestMapping(t *testing.T) {
	query := `
	query {
		users(where: {name: noorain}) @mongo {
			_id
			name
		}
	}
	`
	t.Fatal(parseGraphQLQuery(query))
}
