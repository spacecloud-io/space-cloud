package auth

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
)

//this test generates a tokenstring even if object is empty is this the behaviour we want
func TestCreateToken(t *testing.T) {
	var authCreateToken = []struct {
		testName, wantThis, secretKey string
		IsTokenInvalid                bool
		IsErrExpected                 bool
		object                        map[string]interface{}
	}{
		{testName: "Successful Test", IsTokenInvalid: false, secretKey: "mySecretkey", wantThis: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", object: map[string]interface{}{"id": "internal-sc-user"}},
		{testName: "Test Case-Invalid Token", IsTokenInvalid: true, IsErrExpected: false, secretKey: "mySecretkey", wantThis: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw", object: map[string]interface{}{"id": "internal-scuser"}},
		{testName: "Invalid Test Case-Empty Object", IsTokenInvalid: true, IsErrExpected: false, secretKey: "mySecretkey", wantThis: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw"},
	}
	authModule := Init("1", &crud.Module{}, &schema.Schema{}, false)
	for _, test := range authCreateToken {
		t.Run(test.testName, func(t *testing.T) {
			authModule.SetSecret(test.secretKey)
			tokenString, err := authModule.CreateToken(test.object)
			if (err != nil) != test.IsErrExpected {
				t.Error("Got Error", err, "Wanted Error ", test.IsErrExpected)
			}
			if !reflect.DeepEqual(tokenString, test.wantThis) && !test.IsTokenInvalid {
				t.Error("Got Token", tokenString, "Wanted Token ", test.wantThis)
			}
		})
	}
}

func TestIsTokenInternal(t *testing.T) {
	var authCreateToken = []struct {
		testName, token, secretKey string
		IsErrExpected              bool
	}{
		{testName: "Unsuccessful Test-Token has not been internally created", secretKey: "mySecretkey", IsErrExpected: true, token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc"},
		{testName: "Unsuccessful Test-Signature is Invalid", IsErrExpected: true, secretKey: "mysecretkey", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.MKIZkrXy6nUMu5ejqiYKl7EOU1TxEoKTOww-eoQm6Lw"},
		{testName: "Successful Test Case", IsErrExpected: false, secretKey: "mySecretkey", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImludGVybmFsLXNjLXVzZXIifQ.k3OcidcCnshBOGtzpprfV5Fhl2xWb6sjzPZH3omDDpw"},
	}
	authModule := Init("1", &crud.Module{}, &schema.Schema{}, false)
	for _, test := range authCreateToken {
		t.Run(test.testName, func(t *testing.T) {
			authModule.SetSecret(test.secretKey)
			err := authModule.IsTokenInternal(test.token)
			if (err != nil) != test.IsErrExpected {
				t.Error("Got This ", err, "Wanted Error-", test.IsErrExpected)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	var testCases = []struct {
		name          string
		IsErrExpected bool
		secretKey     string
		token         string
		wantThis      TokenClaims
		reason        error
	}{
		{name: "Test should successfully parse a token", IsErrExpected: false, secretKey: "mySecretkey", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc", wantThis: TokenClaims{"token1": "token1value", "token2": "token2value"}},
		{name: "Test should fail if signing method not HS256", IsErrExpected: true, secretKey: "mySecretkey", token: "eyJhbGciOiJQUzM4NCIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.nakZ1JcYWHcXcG1ZfIY7mJNwcVPQ7U1HvuLEsG9fyz-H9ig3ql8BiI3T-7A2PHe-lBIxjS7hXx8O8lxMg7y7rqUHtPLAGOuCd4Ft88KupgPcF5w-KVpeSgWl598zNLWqJpjcwiPewt3gsU6pwSaTz24JmfZQRrDX8KOtejaGs5OECdk2dDW2rwO98npNX39yYx6eSfZbXCLJ7wIhT3UDbuaOGHnD3wyEtih013NDrnkvVXJRKXUwF7F-g31NWgEgVt-tWkR5vcBBSRYKzIbD7-wxpV4ifLp_XdbVNl3Uf7ja6FeUnGq1Pb9AnAY7lD4Rk7sYQe4P-ATHtkgSg5levw"},
		{name: "Test should fail for an invalid token", IsErrExpected: true, secretKey: "mysecretkey", token: "1234.abcd"},
		{name: "Test should fail for invalid signature-illegal base64 data at input", IsErrExpected: true, secretKey: "mySecretkey", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N"},
		{name: "Test should fail as invalid secret key-invalid signature", IsErrExpected: true, secretKey: "mysecretkey", token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbjEiOiJ0b2tlbjF2YWx1ZSIsInRva2VuMiI6InRva2VuMnZhbHVlIn0.h3jo37fYvnf55A63N-uCyLj9tueFwlGxEGCsf7gCjDc"},
	}

	authModule := Init("1", &crud.Module{}, &schema.Schema{}, false)
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			authModule.SetConfig("default", test.secretKey, config.Crud{}, &config.FileStore{}, &config.ServicesModule{}, &config.Eventing{})
			tokenClaims, err := authModule.parseToken(test.token)
			if (err != nil) != test.IsErrExpected {
				t.Error(test.name, ": Got:", err, "Wanted Error:", test.IsErrExpected)
			}
			if !reflect.DeepEqual(test.wantThis, tokenClaims) {
				t.Error(test.name, ": Got:", tokenClaims, "Want:", test.wantThis, "Reason:", err)
			}

		})
	}
}
