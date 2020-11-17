package admin

import (
	"context"
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

var licenseMode = licenseModeOnline
var licensePublicKey = ""

const maxLicenseFetchErrorCount = 5
const licenseModeOnline = "online"
const licenseModeOffline = "offline"

// Manager manages all admin transactions
type Manager struct {
	lock         sync.RWMutex
	quotas       model.UsageQuotas
	plan         string
	user         *config.AdminUser
	license      *config.License
	integrations config.Integrations

	services model.ScServices
	isProd   bool

	licenseRenewalDate string
	clusterName        string

	syncMan        model.SyncManAdminInterface
	integrationMan IntegrationInterface

	nodeID, clusterID      string
	licenseFetchErrorCount int
	// Config for enterprise
	publicKey *rsa.PublicKey
}

// New creates a new admin manager instance
func New(nodeID, clusterID string, isDev bool, adminUserInfo *config.AdminUser) *Manager {
	m := new(Manager)
	m.nodeID = nodeID
	m.isProd = !isDev // set inverted
	m.clusterID = clusterID
	// Initialise all config
	m.license = new(config.License)
	m.integrations = make(config.Integrations)
	m.user = adminUserInfo
	m.quotas = model.UsageQuotas{MaxDatabases: 1, MaxProjects: 1}

	// Start the background routines
	go m.licenseRenewalCumValidationRoutine()
	go m.fetchPublicKeyRoutine()

	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Starting gateway in %s licensing mode", licenseMode), nil)

	return m
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
func (m *Manager) SetConfig(licenseConfig *config.License) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.license = licenseConfig
	// Create a unique session id if in offline mode
	if licenseMode == licenseModeOffline {
		isLeader, err := m.syncMan.CheckIfLeaderGateway(m.nodeID)
		if err != nil {
			return helpers.Logger.LogError("validate-license-without-lock", "Unable to check who is the current leader gateway", err, nil)
		}
		if isLeader && m.license.LicenseKey == "" || m.license.LicenseValue == "" {
			// Set the licenseKey and value with unique values
			m.license.LicenseKey = ksuid.New().String()
			m.license.LicenseValue = ksuid.New().String()

			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting license key & value", map[string]interface{}{"key": m.license.LicenseKey, "value": m.license.LicenseValue})

			go func() {
				if err := m.syncMan.SetLicense(context.Background(), m.license); err != nil {
					_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to set admin config with session id", nil, nil)
				}
			}()
			helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Successfully set session id", map[string]interface{}{"sessionId": m.license.LicenseKey + m.license.LicenseValue})
			return nil
		}
	}

	if m.isRegistered() {
		if err := m.validateLicenseWithoutLock(m.services); err != nil {
			return err
		}
	} else {
		m.resetQuotasWithoutLock()
	}
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
func (m *Manager) LoadEnv() (bool, string, model.UsageQuotas, string, string, string, string, string, string, string, error) {
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
	sessionID, err := m.GetSessionID()
	if err != nil {
		return false, "", model.UsageQuotas{}, "", "", "", "", "", "", "", err
	}
	return m.isProd, m.plan, m.quotas, loginURL, m.clusterName, m.licenseRenewalDate, m.license.LicenseKey, m.license.LicenseValue, sessionID, licenseMode, nil
}
