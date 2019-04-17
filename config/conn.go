package config

/**
 * @author ollykel, ollykel416@gmail.com
 * @date Apr 20, 2019 
 * Standard configuration struct for SQL databases, including all fields
 * necessary to initialize a connection.
 */

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"encoding/json"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

// setting for whether or not to require ssl, and how to implement
// defaults to "Require"
type sslMode int

const (
	Require sslMode = 0 // default
	Disable sslMode = 1
	VerifyCA sslMode = 1 << 1
	VerifyFull sslMode = 1 << 2
)//-- end sslMode enums

func (mode sslMode) String () string {
	switch mode {
		case Require:
			return "require"
		case Disable:
			return "disable"
		case VerifyCA:
			return "verify-ca"
		case VerifyFull:
			return "verify-full"
		default:
			return "<N/A>"
	}
}//-- end func sslMode.String

func (mode *sslMode) FromString (enum string) error {
	enum = strings.ToLower(enum)
	switch enum {
		case "", "require":
			*mode = Require
		case "disable":
			*mode = Disable
		case "verify-ca", "verifyca":
			*mode = VerifyCA
		case "verify-full", "verifyfull":
			*mode = VerifyFull
		default:
			return fmt.Errorf("unrecognized ssl mode: '%s'", enum)
	}//-- end switch
	return nil
}//-- end func sslMode.FromString

func (mode sslMode) MarshalText () ([]byte, error) {
	return []byte(mode.String()), nil
}//-- end func MarshalYAML

func (mode *sslMode) UnmarshalText (enum []byte) error {
	return mode.FromString(string(enum))
}//-- end func sslMode.UnmarshalYAML

// largely based on go-sql-driver's ConnConfig struct
// 	see: https://godoc.org/go-sql-driver/mysql#ConnConfig
// also based on lib/pq's connection string parameters
// 	see: https://godoc.org/github.com/lib/pq
type ConnConfig struct {
	User		string	`json:"user" yaml:"user"`
	// Auth provides source of password, either env or file
	Auth		string	`json:"auth" yaml:"auth"`
	DBName		string	`json:"dbName" yaml:"dbName"`
	// Protocol ex: "tcp"
	Protocol	string	`json:"protocol" yaml:"protocol"`
	Host		string	`json:"host" yaml:"host"`
	Port		string	`json:"port" yaml:"port"`
	SSLMode		sslMode	`json:"sslMode" yaml:"sslMode"`
	SSLCert		string	`json:"sslCert,omitempty" yaml:"sslCert,omitempty"`
	SSLKey		string	`json:"sslKey,omitempty" yaml:"sslKey,omitempty"`
	SSLRootCert	string	`json:"sslRootCert,omitempty" yaml:"sslRootCert,omitempty"`
	Params		map[string]string `json:"params,omitempty" yaml:"params,omitempty"`
}//-- end SQLConnConfig struct

func (cfg *ConnConfig) GoString () string {
	output := strings.Builder{}
	enc := json.NewEncoder(&output)
	enc.SetIndent("", "  ")
	enc.Encode(cfg)
	return output.String()
}//-- end func ConnConfig.GoString

// ensures that all required fields are provided, throws error otherwise
func (cfg *ConnConfig) validate () error {
	missing := make([]string, 0, 6)//-- 6 = number of req. fields
	if cfg.User == "" { missing = append(missing, "user") }
	if cfg.Auth == "" { missing = append(missing, "auth") }
	if cfg.DBName == "" { missing = append(missing, "dbName") }
	if cfg.Protocol == "" { missing = append(missing, "protocol") }
	if cfg.Host == "" { missing = append(missing, "host") }
	if cfg.Port == "" { missing = append(missing, "port") }
	if len(missing) > 0 {
		return fmt.Errorf("sql config missing fields: %v", missing)
	}
	return nil
}//-- end func ConnConfig.validate

func surveyField (dest *string, fieldName string) error {
	return survey.AskOne(
		&survey.Input{Message: fieldName + ":", Default: *dest},
		dest, survey.Required)
}//-- end func surveyField

func (cfg *ConnConfig) surveySSLMode () (err error) {
	var mode string
	err = survey.AskOne(&survey.Select{
		Message: "sslMode:",
		Options: []string{ "require", "disable",
			"verify-ca", "verify-full" },
		Default: cfg.SSLMode.String()},
		&mode, survey.Required)
	if err != nil { return err }
	return cfg.SSLMode.FromString(mode)
}//-- end func ConnConfig.surveySSLMode

func (cfg *ConnConfig) surveyAuth () (err error) {
	var authType string
	err = survey.AskOne(&survey.Select{
		Message: "auth type:",
		Options: []string{ "FILE", "ENV", "STRING" },
		Default: "FILE"},
		&authType, survey.Required)
	if err != nil { return }
	if authType == "STRING" {
		fmt.Fprint(os.Stderr, "\x1b[31mWarning: " +
			"use of STRING is discouraged.\n" +
			"Consider using FILE or ENV instead\x1b[0m\n")
	}
	err = surveyField(&cfg.Auth, "auth"); if err != nil { return }
	cfg.Auth = authType + ":" + cfg.Auth
	return nil
}//-- end func ConnConfig.surveyAuth

// queries user for input, stores values in ConnConfig
// only covers basic fields (NOT including ssl settings)
// defaults to pre-existing field values
func (cfg *ConnConfig) FromSurvey () (err error) {
	fmt.Print("\x1b[1mGetting Database ConnConfig:\x1b[0m\n")
	// get user
	err = surveyField(&cfg.User, "user"); if err != nil { return }
	// get auth
	err = cfg.surveyAuth(); if err != nil { return }
	// get dbName
	err = surveyField(&cfg.DBName, "dbName"); if err != nil { return }
	// get protocol
	err = surveyField(&cfg.Protocol, "protocol"); if err != nil { return }
	// get host
	err = surveyField(&cfg.Host, "host"); if err != nil { return }
	// get port
	err = surveyField(&cfg.Port, "port"); if err != nil { return }
	// get sslMode
	err = cfg.surveySSLMode()
	return
}//-- end func ConnConfigFromSurvey

func readPasswordFromFile (fname string) (string, error) {
	pass, err := ioutil.ReadFile(fname); if err != nil { return "", err }
	output := strings.Trim(string(pass), "\t\n ")
	return output, nil
}//-- end func readPasswordFromFile

func readPasswordFromEnv (envName string) (string, error) {
	pass, exists := os.LookupEnv(envName)
	if !exists {
		return "", fmt.Errorf("env variable '%s' not set", envName)
	}
	return pass, nil
}//-- end func readPasswordFromEnv

func GetPassword(auth string) (string, error) {
	parts := strings.SplitN(auth, ":", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid auth string '%s'", auth)
	}
	authType, value := strings.ToUpper(parts[0]), parts[1]
	switch authType {
		case "STRING":
			fmt.Fprint(os.Stderr, "\x1b[31mWARNING: " +
				"use of plaintext passwords is discouraged\n" +
				"Consider using FILE or ENV\x1b[0m\n")
			return value, nil
		case "FILE":
			return readPasswordFromFile(value)
		case "ENV":
			return readPasswordFromEnv(value)
		default:
			return "", fmt.Errorf("invalid auth type '%s'", authType)
	}//-- end switch authType
}//-- end func GetPassword

// specifies a function that takes a config and returns a DSN string
type ConnStringParser func (cfg *ConnConfig) string


