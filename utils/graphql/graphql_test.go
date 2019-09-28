package graphql

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
)

func TestMapping(t *testing.T) {
	query := `{
		random(where: {}) @mongo {
			random(where: {key2: "random.key2"}) @mongo {
				key2
			}

			users(where:{name: "random.key1.k1"}) @mongo {
				_id
				email
				name 
				random(where:{key2: "random.key2"}) @mongo {
					key1
				}
			}

			key1 {
				SomeKey:k2
			}

			key2		
		}
	}
	`

	f := functions.Init()
	fConf := &config.Functions{Enabled: true, Broker: "nats", Conn: "nats://localhost:4242", Services: map[string]*config.Service{"default": {Functions: map[string]config.Function{"default": {Rule: &config.Rule{Rule: "allow"}}}}}}
	f.SetConfig(fConf)

	c := crud.Init()
	conf := config.Crud{"mongo": &config.CrudStub{Enabled: true, Conn: "mongodb://localhost:27017", Collections: map[string]*config.TableRule{"default": {Rules: map[string]*config.Rule{"read": {Rule: "allow"}}}}}}
	c.SetConfig(conf)

	a := auth.Init(c, nil)
	a.SetConfig("todo-app", "some-secret", conf, nil, fConf, nil)

	graph := New(a, c, f)
	graph.SetConfig("todo-app")

	te := model.GraphQLRequest{}
	te.Query = query

	var wg sync.WaitGroup
	wg.Add(1)

	graph.ExecGraphQLQuery(&te, "", func(op interface{}, err error) {
		defer wg.Done()

		if err != nil {

			t.Fatal(err)
			return
		}

		data, _ := json.MarshalIndent(op, "", "  ")
		t.Log(string(data))
		return
	})

	wg.Wait()
}
