package admin

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (m *Manager) fetchPublicKeyWithLock() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.fetchPublicKeyWithoutLock()
}

func (m *Manager) licenseRenewalCumValidationRoutine() {
	helpers.Logger.LogInfo("license-renewal-cum-validation-routine", "Starting license renewal/validation routine", nil)
	// Create a random ticker
	min := 6
	max := 24
	for {
		randomInt := rand.Intn(max-min) + min
		t := time.Duration(randomInt) * time.Hour
		select {
		case <-time.After(t):
			// Operate if in enterprise mode
			if m.isEnterpriseMode() {
				isLeader, err := m.syncMan.CheckIfLeaderGateway(m.nodeID)
				if err != nil {
					_ = helpers.Logger.LogError("license-renewal-cum-validation-routine", "Unable to check if i am the leader gateway", err, map[string]interface{}{"nodeId": m.nodeID})
					break
				}
				if isLeader && licenseMode == licenseModeOnline {
					helpers.Logger.LogDebug("license-renewal-cum-validation-routine", "Leader renewing the license", map[string]interface{}{"nodeId": m.nodeID})
					m.lock.Lock()
					if err := m.renewLicenseWithoutLock(false); err != nil {
						_ = helpers.Logger.LogError("license-renewal-cum-validation-routine", "Unable to renew license. Has your subscription expired?", err, nil)
						m.lock.Unlock()
						break
					}
					m.lock.Unlock()
				} else {
					// Check if the license has expired
					helpers.Logger.LogDebug("license-renewal-cum-validation-routine", "Follower validating the license", nil)
					m.lock.Lock()
					m.validationRoutine()
					m.lock.Unlock()
				}
			}
		}
	}
}

func (m *Manager) validationRoutine() {
	helpers.Logger.LogInfo("validation-routine", "Starting license validation routine", nil)
	// Number 6 denotes a total time 30 minutes, with an interval of 5 minutes
	maxRetryCount := 6
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	doesExists := false
	var err error
	currentCount := 0
	for range ticker.C {
		currentCount++
		helpers.Logger.LogInfo("validation-routine", fmt.Sprintf("License validation retry count (%d)", currentCount), nil)
		_, doesExists, err = m.validateSessionID(m.services, m.license.License)
		if err != nil {
			helpers.Logger.LogWarn("validation-routine", "License validation has failed, unable to match license session id with gateway services", map[string]interface{}{"error": err})
			continue
		}
		if currentCount == maxRetryCount || doesExists {
			break
		}
	}
	if !doesExists {
		_ = helpers.Logger.LogError("validation-routine", fmt.Sprintf("Follower gateway is not able to validate the license, resetting global license"), nil, map[string]interface{}{"error": err, "services": m.services, "license": m.license.License})
		m.resetQuotasWithoutLock()
	}
}

func (m *Manager) fetchPublicKeyRoutine() {
	helpers.Logger.LogInfo("fetch-public-key-routine", "Starting license public key fetch routine", nil)
	// Create a new ticker
	ticker := time.NewTicker(4 * 7 * 24 * time.Hour) // fetch public once every 4 weeks
	defer ticker.Stop()

	select {
	case <-ticker.C:
		// Operate if in enterprise mode
		if m.isEnterpriseMode() && licenseMode == licenseModeOnline {
			// Fetch the public key periodically
			if err := m.fetchPublicKeyWithLock(); err != nil {
				_ = helpers.Logger.LogError("fetch-public-key-routine", "Could not fetch public key for license file", err, nil)
				break
			}
		}
	}
}

func (m *Manager) fetchPublicKeyWithoutLock() error {
	// Check if offline licensing mode is used
	if licenseMode == licenseModeOffline {
		// Marshal the public key
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(licensePublicKey))
		if err != nil {
			return helpers.Logger.LogError("fetch-public-key-without-lock", "Unable to parse public key from pem", err, nil)
		}

		// Set the public key
		m.publicKey = publicKey
		return nil
	}

	// Fire the http request
	body := map[string]interface{}{
		"timeout": 10,
	}
	data, _ := json.Marshal(body)
	res, err := http.Post(fmt.Sprintf("https://api.spaceuptech.com/v1/api/spacecloud/services/billing/getPublicKey"), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return helpers.Logger.LogError("fetch-public-key-without-lock", "Unable to contact server to fetch public key", err, nil)
	}

	// Decode the response
	v := new(model.GraphqlFetchPublicKeyResponse)
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return helpers.Logger.LogError("fetch-public-key-without-lock", "Unable to unmarshall result returned by space cloud server", err, nil)
	}

	// Check if valid response was received
	if v.Status != http.StatusOK {
		return helpers.Logger.LogError("fetch-public-key-without-lock", v.Message, errors.New(v.Error), nil)
	}

	// Marshal the public key
	publicKey := new(rsa.PublicKey)
	if err = json.Unmarshal([]byte(v.Result), publicKey); err != nil {
		return helpers.Logger.LogError("fetch-public-key-without-lock", "Unable to unmarshall public key returned by space cloud server to specific format", err, nil)
	}

	// Set the public key
	m.publicKey = publicKey
	return nil
}

func (m *Manager) ValidateLicense(services model.ScServices) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.validateLicenseWithoutLock(services)
}

func (m *Manager) validateSessionID(services model.ScServices, license string) (*model.License, bool, error) {
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return nil, false, helpers.Logger.LogError("validate-session-id", "Unable to fetch public key", err, nil)
		}
	}

	licenseObj, err := m.decryptLicense(license)
	if err != nil {
		m.resetQuotasWithoutLock()
		return nil, false, helpers.Logger.LogError("validate-session-id", "Unable to parse license key, resetting space cloud global license", err, map[string]interface{}{"license": license})
	}

	isFound := false
	if licenseMode == licenseModeOffline {
		if licenseObj.SessionID != m.getOfflineLicenseSessionID() {
			return nil, false, helpers.Logger.LogError("validate-session-id", "Unable to validate license key", errors.New("session id of license file doesn't match with internal session id"), nil)
		}
		isFound = true
	} else {
		for _, service := range services {
			if licenseObj.SessionID == service.ID {
				isFound = true
				break
			}
		}
	}

	return licenseObj, isFound, nil
}

func (m *Manager) validateLicenseWithoutLock(services model.ScServices) error {
	licenseObj, isFound, err := m.validateSessionID(services, m.license.License)
	if err != nil {
		return err
	}
	if isFound {
		helpers.Logger.LogInfo("validate-license-without-lock", "Session id successfully validated, setting internal quotas", nil)
		m.setQuotas(licenseObj)
		return nil
	}

	helpers.Logger.LogInfo("validate-license-without-lock", "Session id validation failed,", map[string]interface{}{"services": services, "license": m.license.License})
	isLeader, err := m.syncMan.CheckIfLeaderGateway(m.nodeID)
	if err != nil {
		return helpers.Logger.LogError("validate-license-without-lock", "Unable to check if i am the leader gateway", err, map[string]interface{}{"nodeId": m.nodeID})
	}
	if isLeader && licenseMode == licenseModeOnline {
		if err := m.renewLicenseWithoutLock(false); err != nil {
			m.resetQuotasWithoutLock()
			return helpers.Logger.LogError("validate-license-without-lock", "Unable to renew license, resetting global license. Check if your subscription has expired?", err, nil)
		}
	}

	return nil
}

func (m *Manager) RenewLicense(force bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	isLeader, err := m.syncMan.CheckIfLeaderGateway(m.nodeID)
	if err != nil {
		return helpers.Logger.LogError("renew-license", "Unable to check if i am the leader gateway", err, map[string]interface{}{"nodeId": m.nodeID, "isForce": force})
	}

	if !isLeader {
		return helpers.Logger.LogError("renew-license", "Unable to renew license", errors.New("only the leader can renew the license"), map[string]interface{}{"nodeId": m.nodeID, "isForce": force})
	}

	// Throw error if licensing mode is set to offline
	if licenseMode == licenseModeOffline {
		return helpers.Logger.LogError("renew-license", "Unable to renew license", errors.New("cannot renew license in offline licensing mode"), nil)
	}

	return m.renewLicenseWithoutLock(force)
}

func (m *Manager) renewLicenseWithoutLock(force bool) error {
	// Marshal the request body
	sessionID := selectRandomSessionID(m.services)
	data, _ := json.Marshal(map[string]interface{}{
		"params": model.RenewLicense{
			LicenseKey:       m.license.LicenseKey,
			LicenseValue:     m.license.LicenseValue,
			License:          m.license.License,
			CurrentSessionID: sessionID,
		},
		"timeout": 10,
	})
	helpers.Logger.LogDebug("renew-license-without-lock", "Renewing admin license", map[string]interface{}{"clusterId": m.license.LicenseKey, "clusterKey": m.license.LicenseValue, "sessionId": sessionID})
	// Fire the request
	res, err := http.Post("https://api.spaceuptech.com/v1/api/spacecloud/services/billing/renewLicense", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return helpers.Logger.LogError("renew-license-without-lock", "Unable to contact server to fetch license file from server", err, nil)
	}
	defer func() { _ = res.Body.Close() }()

	// Decode the response
	data, _ = ioutil.ReadAll(res.Body)

	v := new(model.GraphqlFetchLicenseResponse)
	if err := json.Unmarshal(data, v); err != nil {
		return helpers.Logger.LogError("renew-license-without-lock", "Unable to unmarshal result returned by server", err, nil)
	}

	if res.StatusCode != http.StatusOK {
		return helpers.Logger.LogError("renew-license-without-lock", "Invalid status code received in response", errors.New(v.Error), nil)
	}

	// Check if response is valid
	if v.Status != http.StatusOK {
		m.licenseFetchErrorCount++
		_ = helpers.Logger.LogError("renew-license-without-lock", fmt.Sprintf("Unable to fetch license file. Retry count - %d", m.licenseFetchErrorCount), errors.New(v.Message), nil)
		if m.licenseFetchErrorCount > maxLicenseFetchErrorCount || force {
			helpers.Logger.LogInfo("renew-license-without-lock", "Max retry limit hit.", nil)
			m.ResetQuotas()
			return fmt.Errorf("%s-%s", v.Message, v.Error)
		}
		return nil
	} else {
		m.licenseFetchErrorCount = 0
	}

	licenseObj, isSessionValid, err := m.validateSessionID(m.services, v.Result.License)
	if err != nil {
		return err
	}
	if !isSessionValid {
		return helpers.Logger.LogError("renew-license-without-lock", "Found invalid session id in the newly renewed license", nil, map[string]interface{}{"services": m.services, "license": v.Result.License})
	}

	m.license.License = v.Result.License
	m.setQuotas(licenseObj)

	go func() { _ = m.syncMan.SetLicense(context.TODO(), m.license) }()
	return nil
}

func (m *Manager) ResetQuotas() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.resetQuotasWithoutLock()
}

func (m *Manager) resetQuotasWithoutLock() {
	helpers.Logger.LogInfo("reset-quotas-without-lock", "Resetting space cloud to run in open source mode. You will have to re-register the cluster again.", nil)
	m.quotas.MaxProjects = 1
	m.quotas.MaxDatabases = 1
	m.quotas.IntegrationLevel = 0
	m.plan = "space-cloud-open--monthly"

	if licenseMode == licenseModeOnline {
		m.license.LicenseKey = ""
		m.license.LicenseValue = ""
	}

	m.license.License = ""

	m.clusterName = ""

	isLeader, err := m.syncMan.CheckIfLeaderGateway(m.nodeID)
	if err != nil {
		_ = helpers.Logger.LogError("reset-quotas-without-lock", "Unable to check if i am the leader gateway", err, map[string]interface{}{"nodeId": m.nodeID})
	}
	if isLeader {
		go func() {
			if err := m.syncMan.SetLicense(context.Background(), m.license); err != nil {
				_ = helpers.Logger.LogError("reset-quotas-without-lock", "Unable to save empty license in global config", err, map[string]interface{}{"nodeId": m.nodeID})
			}
		}()
	}
}

func (m *Manager) setQuotas(licenseObj *model.License) {
	// set quotas
	m.quotas.MaxProjects = licenseObj.Meta.ProductMeta.MaxProjects
	m.quotas.MaxDatabases = licenseObj.Meta.ProductMeta.MaxDatabases
	m.quotas.IntegrationLevel = licenseObj.Meta.ProductMeta.IntegrationLevel
	m.clusterName = licenseObj.Meta.LicenseKeyMeta.ClusterName
	m.licenseRenewalDate = licenseObj.LicenseRenewal
	m.plan = licenseObj.Plan
	helpers.Logger.LogInfo("set-quotas", fmt.Sprintf("Gateway is running with %s plan ", licenseObj.Plan), nil)
}

func (m *Manager) isEnterpriseMode() bool {
	return m.isRegistered() && !strings.HasPrefix(m.plan, "space-cloud-open")
}

func (m *Manager) isRegistered() bool {
	return m.license.LicenseKey != "" && m.license.LicenseValue != "" && m.license.License != ""
}

func (m *Manager) decryptLicense(license string) (*model.License, error) {
	obj, err := m.parseLicenseToken(license)
	if err != nil {
		return nil, err
	}

	v := new(model.License)
	if err := mapstructure.Decode(obj, v); err != nil {
		return nil, helpers.Logger.LogError("decrypt-license", "Unable to decode license object to license struct", err, nil)
	}
	return v, nil

}

func (m *Manager) parseLicenseToken(tokenString string) (map[string]interface{}, error) {
	licenseObj, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, helpers.Logger.LogError("parse-license-token", "Unable to parse license", errors.New("invalid signing method type"), nil)
		}

		return m.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := licenseObj.Claims.(jwt.MapClaims); ok && licenseObj.Valid {
		if err := claims.Valid(); err != nil {
			return nil, helpers.Logger.LogError("parse-license-token", "Unable to validate license claims", err, nil)
		}
		return claims, nil
	}

	return nil, helpers.Logger.LogError("parse-license-token", "Unable to parse license token", nil, nil)
}

func selectRandomSessionID(gateways model.ScServices) string {
	if len(gateways) == 0 {
		helpers.Logger.LogWarn("select-random-session-id", fmt.Sprintf("Length of gateways is zero"), nil)
		return ""
	}
	min := 0
	max := len(gateways)
	// get an int from min...max-1 range
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(max-min) + min
	helpers.Logger.LogInfo("select-random-session-id", fmt.Sprintf("Selecting session id (%s)", gateways[index].ID), nil)
	return gateways[index].ID
}

func (m *Manager) getOfflineLicenseSessionID() string {
	return m.license.LicenseKey + m.license.LicenseValue
}
