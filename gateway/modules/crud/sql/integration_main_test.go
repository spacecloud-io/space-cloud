// +build integration

package sql

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

var dbType = flag.String("db_type", "", "db_type of test case to be run")
var connection = flag.String("conn", "", "connection string of the database")

func TestMain(m *testing.M) {
	var customerTable = `type customers {
							id: ID! @primary
							name: String!
							age: Integer!
							height: Float
							is_prime: Boolean! 
							birth_date: DateTime!
							address: JSON!
						}`
	var companiesTable = `type companies {
 						 id: ID! @primary
 						 parent : ID!
 						 name : String!
 						 established_date : DateTime!
 						 kind : Integer!
 						 volume : Float!
 						 is_public : Boolean!
 						 description : JSON!
					}`
	var ordersTable = `type orders {
						id: ID! @primary
						order_date: DateTime!
						amount: Integer!
						is_prime: Boolean,
						product_id: String!
						address: JSON!
						stars: Float!
					}`
	var rawBatch = `type raw_batch {
						id: ID! @primary
						score : Integer!
					}`
	var rawQuery = `type raw_batch {
						id: ID! @primary
						score : Integer!
					}`
	// create sc project
	var projectInfo = `{"name":"myproject","id":"myproject","secrets":[{"secret":"27f6a16bf7864c319e01b7511737407d","isPrimary":true}],"aesKey":"MWJkOTE5ZjVmMGRjNGZiMjg4MDQ0NjQ5MDE0ZWM2MDQ=","contextTime":5,"modules":{"db":{},"eventing":{},"userMan":{},"remoteServices":{"externalServices":{}},"fileStore":{"enabled":false,"rules":[]}}}`
	res, err := http.Post("http://localhost:4122/v1/config/projects/myproject", "application/json", bytes.NewBuffer([]byte(projectInfo)))
	if err != nil {
		log.Printf("Integration test couldn't create companies table - %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Integration test couldn't create companies table got status code %v", res.StatusCode)
		return
	}

	// connect to database
	data, _ := json.Marshal(map[string]interface{}{"enabled": true, "conn": *connection, "type": *dbType, "name": "myproject"})
	res, err = http.Post(fmt.Sprintf("http://localhost:4122/v1/config/projects/myproject/database/%s/config/database-config", *dbType), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Integration test couldn't create companies table - %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Integration test couldn't create companies table got status code %v", res.StatusCode)
		return
	}

	// create table from space cloud
	data, _ = json.Marshal(map[string]interface{}{
		"collections": map[string]interface{}{
			"default": map[string]interface{}{},
			"companies": map[string]interface{}{
				"schema": companiesTable,
			},
			"orders": map[string]interface{}{
				"schema": ordersTable,
			},
			"raw_batch": map[string]interface{}{
				"schema": rawBatch,
			},
			"raw_query": map[string]interface{}{
				"schema": rawQuery,
			},
			"customers": map[string]interface{}{
				"schema": customerTable,
			},
		},
	})
	res, err = http.Post(fmt.Sprintf("http://localhost:4122/v1/config/projects/myproject/database/%s/schema/mutate", *dbType), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Integration test couldn't create companies table - %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Integration test couldn't create companies table got status code %v", res.StatusCode)
		return
	}
	exitVal := m.Run()
	// cleaning
	os.Exit(exitVal)
}
