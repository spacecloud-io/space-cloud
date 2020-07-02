// +build integration

package mgo

import (
	"bytes"
	"context"
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
	flag.Parse()
	//var customerTable = `type customers {
	//						id: ID! @primary
	//						name: String!
	//						age: Integer!
	//						height: Float
	//						is_prime: Boolean!
	//						birth_date: DateTime!
	//						address: JSON!
	//					}`
	//var companiesTable = `type companies {
	//					 id: ID! @primary
	//					 parent : ID!
	//					 name : String!
	//					 established_date : DateTime!
	//					 kind : Integer!
	//					 volume : Float!
	//					 is_public : Boolean!
	//					 description : JSON!
	//				}`
	//var ordersTable = `type orders {
	//					id: ID! @primary
	//					order_date: DateTime!
	//					amount: Integer!
	//					is_prime: Boolean,
	//					product_id: String!
	//					address: JSON!
	//					stars: Float!
	//				}`
	//var rawBatch = `type raw_batch {
	//					id: ID! @primary
	//					score : Integer!
	//				}`
	//var rawQuery = `type raw_query {
	//					id: ID! @primary
	//					score : Integer!
	//				}`

	// create sc project
	var projectInfo = `{"name":"myproject","id":"myproject","secrets":[{"secret":"27f6a16bf7864c319e01b7511737407d","isPrimary":true}],"aesKey":"MWJkOTE5ZjVmMGRjNGZiMjg4MDQ0NjQ5MDE0ZWM2MDQ=","contextTime":5,"modules":{"db":{},"eventing":{},"userMan":{},"remoteServices":{"externalServices":{}},"fileStore":{"enabled":false,"rules":[]}}}`
	res, err := http.Post("http://localhost:4122/v1/config/projects/myproject", "application/json", bytes.NewBuffer([]byte(projectInfo)))
	if err != nil {
		log.Printf("Integration test couldn't create project (myproject) in space cloud - %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Integration test couldn't create project (myproject) in space cloud got status code %v", res.StatusCode)
		return
	}

	// connect to database
	data, _ := json.Marshal(map[string]interface{}{"enabled": true, "conn": *connection, "type": *dbType, "name": "myproject"})
	res, err = http.Post(fmt.Sprintf("http://localhost:4122/v1/config/projects/myproject/database/%s/config/database-config", *dbType), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Integration test couldn't add database in space cloud table - %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		v := map[string]interface{}{}
		json.NewDecoder(res.Body).Decode(&v)
		log.Printf("Integration test couldn't add database in space cloud got status code %v -error %v", res.StatusCode, v["error"])
		return
	}

	log.Printf("Running tests")
	exitVal := m.Run()

	// cleaning
	log.Printf("Cleaning data")

	// delete database config from space cloude
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:4122/v1/config/projects/myproject/database/%s/config/database-config", *dbType), nil)
	if err != nil {
		log.Printf("Integration test couldn't call space-cloud from removing database config - %v", err)
		return
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Integration test couldn't remove database config from space-cloud - %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Integration test couldn't remove database config from space-cloud got status code %v", res.StatusCode)
		return
	}

	db, err := Init(true, *connection, "myproject")
	if err != nil {
		log.Println("Create() Couldn't establishing connection with database", dbType)
		return
	}
	// clear data
	if err := db.client.Database("myproject").Drop(context.Background()); err != nil {
		log.Println("Create() Couldn't truncate table", err)
	}

	os.Exit(exitVal)
}
