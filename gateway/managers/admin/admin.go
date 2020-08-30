package admin

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

var licenseMode = "online"
var licensePublicKey = ""

const maxLicenseFetchErrorCount = 5

// Manager manages all admin transactions
type Manager struct {
	lock   sync.RWMutex
	config *config.Admin
	quotas model.UsageQuotas
	plan   string
	user   *config.AdminUser
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
	m.config = new(config.Admin)
	m.user = adminUserInfo
	m.quotas = model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1}

	// Start the background routines
	go m.licenseRenewalRoutine()
	go m.fetchPublicKeyRoutine()

	utils.LogInfo(fmt.Sprintf("Starting gateway in %s licensing mode", licenseMode), "admin", "new")

	return m
}

func (m *Manager) startOperation(license string, isInitialCall bool) error {
	logrus.Infoln("Starting gateway in enterprise mode")

	// Fetch the public key if it does't already exist
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return utils.LogError("Unable to fetch public key", "admin", "startOperation", err)
		}
	}

	// Parse the license
	licenseObj, err := m.decryptLicense(license)
	if err != nil {
		return utils.LogError("Unable to decrypt license key", "admin", "startOperation", err)
	}

	// We have a problem if our session id does not match with the license's session id
	if m.sessionID != licenseObj.SessionID {

		// There cannot be a mismatch unless the gateway just started while being in online mode. For anytime else, throw an error.
		if !isInitialCall {

			// Reset quotas and admin config to defaults
			_ = utils.LogError("Invalid license file provided. Did you change the license key yourself?", "admin", "startOperation", errors.New("session id mismatch while setting admin config"))
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
func (m *Manager) SetConfig(config *config.Admin, isInitialCall bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Set the admin config
	m.config = config

	// Create a unique session id if in offline mode
	if licenseMode == "offline" {
		if m.config.LicenseKey == "" || m.config.LicenseValue == "" {
			// Set the licenseKey and value with unique values
			m.config.LicenseKey = ksuid.New().String()
			m.config.LicenseValue = ksuid.New().String()

			utils.LogDebug("Setting session id", "admin", "set-config", map[string]interface{}{"key": m.config.LicenseKey, "value": m.config.LicenseValue})

			go func() {
				if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
					_ = utils.LogError("Unable to set admin config with session id", "admin", "set-config", nil)
				}
			}()
			return nil
		}

		m.sessionID = m.config.LicenseKey + m.config.LicenseValue
		utils.LogDebug("Successfully set session id", "admin", "set-config", map[string]interface{}{"sessionId": m.sessionID})
	}

	// Check if the cluster is registered
	if m.isRegistered() {
		if m.checkIfLeaderGateway() && licenseMode == "online" {
			// Only the leader gateway can handle licensing information
			return m.startOperation(config.License, isInitialCall)
		} else {
			return m.setQuotas(config.License)
		}
	}

	utils.LogInfo("Gateway running in open source mode", "admin", "SetConfig")
	// Reset quotas defaults
	m.quotas.MaxProjects = 1
	m.quotas.MaxDatabases = 1
	m.quotas.IntegrationLevel = 0
	m.plan = "space-cloud-open--monthly"
	return nil
}

// GetConfig returns the admin config
func (m *Manager) GetConfig() *config.Admin {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.config
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

	return m.isProd, m.plan, m.quotas, loginURL, m.clusterName, m.licenseRenewalDate, m.config.LicenseKey, m.config.LicenseValue, m.sessionID, licenseMode
}
