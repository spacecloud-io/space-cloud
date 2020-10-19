package admin

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

var licenseMode = "online"
var licensePublicKey = ""

const maxLicenseFetchErrorCount = 5

// Manager manages all admin transactions
type Manager struct {
	lock         sync.RWMutex
	quotas       model.UsageQuotas
	plan         string
	user         *config.AdminUser
	license      *config.License
	integrations config.Integrations

	isProd bool

	licenseRenewalDate string
	clusterName        string

	syncMan        model.SyncManAdminInterface
	integrationMan IntegrationInterface

	nodeID, clusterID      string
	licenseFetchErrorCount int
	// Config for enterprise
	sessionID string
	publicKey *rsa.PublicKey
}

// New creates a new admin manager instance
func New(nodeID, clusterID string, isDev bool, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.nodeID = nodeID
	m.isProd = !isDev // set inverted
	m.clusterID = clusterID
	if m.checkIfLeaderGateway() {
		m.sessionID = ksuid.New().String()
	}
	// Initialise all config
	m.license = new(config.License)
	m.integrations = make(config.Integrations)
	m.user = adminUserInfo
	m.quotas = model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1}

	// Start the background routines
	go m.licenseRenewalRoutine()
	go m.fetchPublicKeyRoutine()

	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Starting gateway in %s licensing mode", licenseMode), nil)

	return m
}

func (m *Manager) startOperation(license string, isInitialCall bool) error {
	helpers.Logger.LogInfo("", "Starting gateway in enterprise mode", nil)

	// Fetch the public key if it does't already exist
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to fetch public key", err, nil)
		}
	}

	// Parse the license
	licenseObj, err := m.decryptLicense(license)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to decrypt license key", err, nil)
	}

	// We have a problem if our session id does not match with the license's session id
	if m.sessionID != licenseObj.SessionID {

		// There cannot be a mismatch unless the gateway just started while being in online mode. For anytime else, throw an error.
		if !isInitialCall {

			// Reset quotas and admin config to defaults
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid license file provided. Did you change the license key yourself?", errors.New("session id mismatch while setting admin config"), nil)
			m.ResetQuotas()
			return nil
		}

		// Renew the license to update the license session id with the current id
		if err := m.renewLicenseWithoutLock(true); err != nil {
			return err
		}
		return nil
	}

	// set quotas
	m.quotas.MaxProjects = licenseObj.Meta.ProductMeta.MaxProjects
	m.quotas.MaxDatabases = licenseObj.Meta.ProductMeta.MaxDatabases
	m.quotas.IntegrationLevel = licenseObj.Meta.ProductMeta.IntegrationLevel
	m.clusterName = licenseObj.Meta.LicenseKeyMeta.ClusterName
	m.licenseRenewalDate = licenseObj.LicenseRenewal
	m.plan = licenseObj.Plan

	return nil
}

func (m *Manager) SetSyncMan(s model.SyncManAdminInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.syncMan = s
}

func (m *Manager) SetIntegrationMan(i IntegrationInterface) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.integrationMan = i
}

// SetConfig sets the admin config
func (m *Manager) SetConfig(licenseConfig *config.License, isInitialCall bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Set the admin config
	m.license = licenseConfig

	// Create a unique session id if in offline mode
	if licenseMode == "offline" {
		if m.license.LicenseKey == "" || m.license.LicenseValue == "" {
			// Set the licenseKey and value with unique values
			m.license.LicenseKey = ksuid.New().String()
			m.license.LicenseValue = ksuid.New().String()

			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting session id", map[string]interface{}{"key": m.license.LicenseKey, "value": m.license.LicenseValue})

			go func() {
				if err := m.syncMan.SetLicense(context.Background(), m.license); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set admin config with session id", nil, nil)
				}
			}()
			return nil
		}

		m.sessionID = m.license.LicenseKey + m.license.LicenseValue
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Successfully set session id", map[string]interface{}{"sessionId": m.sessionID})
	}

	// Check if the cluster is registered
	if m.isRegistered() {
		if m.checkIfLeaderGateway() && licenseMode == "online" {
			// Only the leader gateway can handle licensing information
			return m.startOperation(licenseConfig.License, isInitialCall)
		} else {
			return m.setQuotas(licenseConfig.License)
		}
	}

	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), "Gateway running in open source mode", nil)
	// Reset quotas defaults
	m.quotas.MaxProjects = 3
	m.quotas.MaxDatabases = 3
	m.quotas.IntegrationLevel = 10
	m.plan = "space-cloud-open--monthly"
	return nil
}

// GetConfig returns the admin config
func (m *Manager) GetConfig() *config.License {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.license
}

func (m *Manager) SetIntegrationConfig(integrations config.Integrations) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.integrations = integrations
}

// LoadEnv gets the env
func (m *Manager) LoadEnv() (bool, string, model.UsageQuotas, string, string, string, string, string, string, string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loginURL := "/mission-control/login"

	// Invoke integration hooks
	hookResponse := m.integrationMan.InvokeHook(ctx, model.RequestParams{
		Resource: "load-env",
		Op:       "read",
	})
	if hookResponse.CheckResponse() {
		if err := hookResponse.Error(); err == nil {
			loginURL = hookResponse.Result().(map[string]interface{})["loginUrl"].(string)
		}
	}

	return m.isProd, m.plan, m.quotas, loginURL, m.clusterName, m.licenseRenewalDate, m.license.LicenseKey, m.license.LicenseValue, m.sessionID, licenseMode
}
