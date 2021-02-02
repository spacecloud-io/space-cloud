package auth

import (
	"context"
	"encoding/base64"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SetProjectConfig set project config of auth module
func (m *Module) SetProjectConfig(projectConfig *config.ProjectConfig) error {
	m.Lock()
	defer m.Unlock()

	m.project = projectConfig.ID
	if projectConfig.SecretSource == "admin" {
		projectConfig.Secrets = []*config.Secret{{KID: utils.AdminSecretKID, Secret: m.adminMan.GetSecret(), IsPrimary: true, Alg: config.HS256}}
	}

	if err := m.jwt.SetSecrets(projectConfig.Secrets); err != nil {
		return err
	}

	decodedAESKey, err := base64.StdEncoding.DecodeString(projectConfig.AESKey)
	if err != nil {
		return err
	}
	m.aesKey = decodedAESKey
	return nil
}

// SetConfig set the rules and secret key required by the auth block
func (m *Module) SetConfig(ctx context.Context, fileStoreType string, projectConfig *config.ProjectConfig, dbRules config.DatabaseRules, dbPreparedRules config.DatabasePreparedQueries, fileStoreRules config.FileStoreRules, remoteServices config.Services, eventingRules config.EventingRules, securityFunctions config.SecurityFunctions) error {

	if err := m.SetProjectConfig(projectConfig); err != nil {
		return err
	}

	m.SetDatabaseRules(dbRules)
	m.SetDatabasePreparedQueryRules(dbPreparedRules)
	m.SetFileStoreRules(fileStoreRules)
	m.SetEventingRules(eventingRules)
	m.SetRemoteServiceConfig(remoteServices)
	m.SetSecurityFunctionConfig(securityFunctions)

	return nil
}

// CloseConfig closes go routines and initializes maps
func (m *Module) CloseConfig() {
	m.Lock()
	defer m.Unlock()

	m.jwt.Close()
	m.funcRules = map[string]*config.Service{}
	m.eventingRules = map[string]*config.Rule{}
	m.fileRules = []*config.FileRule{}
	m.dbRules = map[string]*config.DatabaseRule{}
}

// SetRemoteServiceConfig sets the service module config
func (m *Module) SetRemoteServiceConfig(remoteServices config.Services) {
	m.Lock()
	defer m.Unlock()
	m.funcRules = remoteServices
}

// SetSecurityFunctionConfig sets the security function  config
func (m *Module) SetSecurityFunctionConfig(securityFunctions config.SecurityFunctions) {
	m.Lock()
	defer m.Unlock()
	m.securityFunctions = securityFunctions
}

// SetFileStoreRules sets the file store module config
func (m *Module) SetFileStoreRules(fileRules config.FileStoreRules) {
	m.Lock()
	defer m.Unlock()
	if fileRules == nil {
		return
	}
	temp := make([]*config.FileRule, 0)
	for _, rule := range fileRules {
		temp = append(temp, rule)
	}
	sortFileRule(temp)
	m.fileRules = temp
}

// SetFileStoreType sets file story type
func (m *Module) SetFileStoreType(fileStoreType string) {
	m.Lock()
	defer m.Unlock()
	m.fileStoreType = fileStoreType
}

// SetEventingRules sets the eventing config
func (m *Module) SetEventingRules(eventingRules config.EventingRules) {
	m.Lock()
	defer m.Unlock()
	m.eventingRules = eventingRules
}

// SetDatabaseRules sets the crud module config
func (m *Module) SetDatabaseRules(dbRules config.DatabaseRules) {
	m.Lock()
	defer m.Unlock()
	m.dbRules = dbRules
}

// SetDatabasePreparedQueryRules set prepared query rules of auth module
func (m *Module) SetDatabasePreparedQueryRules(dbPreparedRules config.DatabasePreparedQueries) {
	m.Lock()
	defer m.Unlock()
	m.dbPrepQueryRules = dbPreparedRules
}

// SetMakeHTTPRequest sets the http request
func (m *Module) SetMakeHTTPRequest(function utils.TypeMakeHTTPRequest) {
	m.Lock()
	defer m.Unlock()

	m.makeHTTPRequest = function
}
