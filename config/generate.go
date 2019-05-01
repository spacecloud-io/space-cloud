package config

import (
	"fmt"
	"os"
	"log"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/spaceuptech/space-cloud/utils"
)

type input struct {
	Conn		string
	ConnConfig	*ConnConfig
	PrimaryDB	string
	ID			string
}

/**
 * @author ollykel, ollykel416@gmail.com
 * @date Apr 19, 2019
 * Formats db into appropriate name for config
 */
func formatDBName (option string) (dbName string) {
	switch option {
		case "mongo", "mongodb":
			return "mongo"
		case "mysql", "postgres":
			return "sql-" + option
		default:
			return ""
	}//-- end switch
}//-- end func formatDBName

func getProjectId (dest *input) (err error) {
	dir, err := os.Getwd()
	if err != nil { log.Panic(err.Error()) }
	array := strings.Split(dir, "/")
	dir = array[len(array)-1]
	err = survey.AskOne(
		&survey.Input{Message: "project name:",
		Default: formatProjectID(dir)},
		&dest.ID, survey.Required)
	if err != nil { return err }
	dest.ID = formatProjectID(dest.ID)
	return err
}//-- end func getProjectId

func getDB (dest *input) (err error) {
	return survey.AskOne(&survey.Select{
		Message: "primary database:",
		Options: []string{"mongo", "mysql", "postgres"},
		Default: "mongo",
	}, &dest.PrimaryDB, survey.Required)
}//-- end func getDB

func getConnString (dest *input) (err error) {
	return survey.AskOne(&survey.Input{
		Message: "connection string (" + dest.PrimaryDB + "}",
		Default: getConnectionString(dest.PrimaryDB)},
		&dest.Conn, survey.Required)
}//-- end func getConnString

/**
 * @author ollykel, ollykel416@gmail.com
 * @date Apr 22, 2019
 * Surveys user for fields for sql.Config, which will be used to generate
 * the connection string at runtime
 */
func getConnectionConfig (dest *input) (err error) {
	// config with reasonable defaults
	dest.ConnConfig = &ConnConfig{
		User: "user",
		Auth: "password",
		DBName: dest.ID,
		Protocol: "tcp", Host: "localhost" }
	switch dest.PrimaryDB {
		case string(utils.Mongo):
			dest.ConnConfig.Port = "27017"
		case string(utils.MySQL):
			// set defaults for MySQL
			dest.ConnConfig.Port = "3306"
		case string(utils.Postgres):
			// set defaults for Postgres
			dest.ConnConfig.Port = "5432"
		default:
			return fmt.Errorf("unrecognized database '%s'", dest.PrimaryDB)
	}//-- end switch dest.PrimaryDB
	return dest.ConnConfig.FromSurvey()
}//-- end func getConnectionConfig

const generateConfigPrompt = `
This utility walks you through creating a config.yaml file for your
space-cloud project.
It only covers the most essential configurations and suggests
sensible defaults.

Press ^C at any time to quit.
`//-- end const generateConfigPrompt

// starts the interactive cli to generate config file
func GenerateConfig() (err error) {
	fmt.Print(generateConfigPrompt)

	in := input{}

	// Ask for the project id
	if err = getProjectId(&in); err != nil { return err }

	// Ask the primary db
	if err = getDB(&in); err != nil { return err }

	// prepend sql- for mysql, etc
	in.PrimaryDB = formatDBName(in.PrimaryDB)

	// Ask for db connection info
	in.Conn = ""//-- remove default connection string
	if err = getConnectionConfig(&in); err != nil { return err }

	return writeConfig(&in)
}//-- end func GenerateConfig

func writeConfig(in *input) (err error) {
	configReader := strings.NewReader(defaultTemplate)
	proj, err := LoadConfig(configReader, "yaml"); if err != nil { return }

	proj.ID = in.ID

	crud := proj.Modules.Crud
	crud[in.PrimaryDB] = crud["primary"]
	delete(crud, "primary")

	crud[in.PrimaryDB].Conn = in.Conn
	crud[in.PrimaryDB].Connection = in.ConnConfig

	return proj.Save()
}//-- end func writeConfig

func formatProjectID(id string) string {
	return strings.Join(strings.Split(strings.ToLower(id), " "), "-")
}

func getConnectionString(db string) string {
	switch db {
	case string(utils.Mongo):
		return "mongodb://localhost:27017"
	case string(utils.MySQL):
		return "user:my-secret-pw@/test"
	case string(utils.Postgres):
		return "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"
	default:
		return "localhost"
	}
}
