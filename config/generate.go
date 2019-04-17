package config

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/spaceuptech/space-cloud/utils"
)

type input struct {
	Conn      string
	PrimaryDB string
	ID        string
}

// GenerateConfig started the interactive cli to generate config file
func GenerateConfig() error {
	fmt.Println()
	fmt.Println("This utility walks you through creating a config.yaml file for your space-cloud project.")
	fmt.Println("It only covers the most essential configurations and suggests sensible defaults.")
	fmt.Println()
	fmt.Println("Press ^C any time to quit.")

	i := new(input)

	// Ask the project id
	dir, _ := os.Getwd()
	array := strings.Split(dir, string(os.PathSeparator))
	dir = array[len(array)-1]
	err := survey.AskOne(&survey.Input{Message: "project name:", Default: formatProjectID(dir)}, &i.ID, survey.Required)
	if err != nil {
		return err
	}
	i.ID = formatProjectID(i.ID)

	// Ask the primary db
	err = survey.AskOne(&survey.Select{
		Message: "primary database:",
		Options: []string{"mongo", "mysql", "postgres"},
		Default: "mongo",
	}, &i.PrimaryDB, survey.Required)
	if err != nil {
		return err
	}
	if i.PrimaryDB == "mysql" || i.PrimaryDB == "postgres" {
		i.PrimaryDB = "sql-" + i.PrimaryDB
	}

	// Ask for the connection string
	err = survey.AskOne(&survey.Input{Message: "connection string (" + i.PrimaryDB + ")", Default: getConnectionString(i.PrimaryDB)}, &i.Conn, survey.Required)
	if err != nil {
		return err
	}

	return writeConfig(i)
}

func writeConfig(i *input) error {
	f, err := os.Create("./" + i.ID + ".yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("config").Parse(templateString)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, i)
	if err != nil {
		return err
	}

	return f.Sync()
}

func formatProjectID(id string) string {
	return strings.Join(strings.Split(strings.ToLower(id), " "), "-")
}

func getConnectionString(db string) string {
	switch db {
	case string(utils.Mongo):
		return "mongodb://localhost:27017"
	case string(utils.MySQL):
		return "root:my-secret-pw@/test"
	case string(utils.Postgres):
		return "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"
	default:
		return "localhost"
	}
}
