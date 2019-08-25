package crud

import (
	"testing"

	"github.com/spaceuptech/space-cloud/config"
)

func TestParseSchema(t *testing.T) {
	query := `
	type Tweet {
		id: ID! @id
		createdAt: DateTime! @createdAt
		text: String
		owner: [User] @relation(link: INLINE)
		location: Location!
	  }
	  
	  type User {
		id: ID! @id
		createdAt: DateTime! @createdAt
		updatedAt: DateTime! @updatedAt
		handle: String! @unique
		name: String
		tweets: [Tweet!]!
	  }
	  
	  type Location {
		id: ID! @id
		latitude: Float!
		longitude: Float!
	  }
	`
	v := config.Crud{
		"mongo": &config.CrudStub{
			Collections: map[string]*config.TableRule{
				"tweet": &config.TableRule{
					Schema: query,
				},
			},
		},
	}
	m := Init()

	t.Run("Schema Parser", func(t *testing.T) {
		output, err := m.ParseSchema(v)
		if err != nil {
			t.Fatal(err)
		}
		// b, err := json.MarshalIndent(output, "", "  ")
		// if err != nil {
		// 	fmt.Println("error:", err)
		// }
		// fmt.Print(string(b))
		t.Log("Logging Test Output :: ", output)
	})
}
