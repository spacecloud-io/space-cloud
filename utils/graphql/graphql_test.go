package graphql

import (
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
)

func TestMapping(t *testing.T) {
	query := `
	query {
		random(where: $find, skip: 87) @mongo(col: random) {
			key1 {
				SomeKey:k2
			}
			key2
			
			users(where:{name: random__key1__k1}) @mongo {
				_id
				email
				
			}
		}
	}
	mutation {
		user(where: {}, set: {}, op: upsert)
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

	output, err := graph.ExecGraphQLQuery(query)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
}
