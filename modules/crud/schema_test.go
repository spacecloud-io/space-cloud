package crud

import (
	"testing"

	"github.com/spaceuptech/space-cloud/model"

	"github.com/spaceuptech/space-cloud/config"
)

var query = `
type Tweet {
	id: ID! @id
	createdAt: DateTime! @createdAt
	text: String
	owner: [Integer!]! @relation(link: INLINE)
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
var v = config.Crud{
	"mongo": &config.CrudStub{
		Collections: map[string]*config.TableRule{
			"tweet": &config.TableRule{
				Schema: query,
			},
			"location": &config.TableRule{
				Schema: query,
			},
			"user": &config.TableRule{
				Schema: query,
			},
		},
	},
}

func TestParseSchema(t *testing.T) {

	m := Init()

	t.Run("Schema Parser", func(t *testing.T) {
		err := m.ParseSchema(v)
		if err != nil {
			t.Fatal(err)
		}
		// b, err := json.MarshalIndent(m.schema, "", "  ")
		// if err != nil {
		// 	fmt.Println("error:", err)
		// }
		// fmt.Print(string(b))
		t.Log("Logging Test Output :: ", m.schema)
	})
}

func TestValidateSchema(t *testing.T) {

	var arr []interface{}
	// str := []string{"sharad", "regoti", "atharva"}
	str := []int{1, 2, 3}
	for _, v := range str {
		arr = append(arr, v)
	}

	req := model.CreateRequest{
		Document: []map[string]interface{}{
			{
				"id":        "dfdsairfa",
				"createdAt": 986413662654,
				"text":      "Hello World!",
				"location": map[string]interface{}{
					"id":        "locatoinid",
					"latitude":  5.6,
					"longitude": 312.3,
				},
				"owner": arr,
			},
		},
	}

	tdd := []struct {
		dbName, coll, description string
		value                     model.CreateRequest
	}{{
		dbName:      "mongo",
		coll:        "tweet",
		description: "checking User defined type",
		value:       req,
	}}

	m := Init()
	err := m.ParseSchema(v)
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range tdd {
		t.Run(val.description, func(t *testing.T) {
			err := m.ValidateSchema(val.dbName, val.coll, &val.value)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}
