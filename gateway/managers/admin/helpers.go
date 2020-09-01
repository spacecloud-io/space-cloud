package admin

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func (m *Manager) licenseRenewalRoutine() {
	// Create a new ticker
	ticker := time.NewTicker(24 * time.Hour) // renew license every day
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Operate if in enterprise mode
			if m.isEnterpriseMode() {
				if m.checkIfLeaderGateway() && licenseMode == "online" {
					// Fetch the public key periodically
					if err := m.RenewLicense(false); err != nil {
						_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to renew license. Has your subscription expired?", err, nil)
						break
					}
					go func() {
						if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
							_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to save admin config", err, nil)
						}
					}()
				} else {
					// Check if the license has expired
					_ = m.ValidateLicense()
				}
			}
		}
	}
}

func (m *Manager) fetchPublicKeyRoutine() {
	// Create a new ticker
	ticker := time.NewTicker(4 * 7 * 24 * time.Hour) // fetch public once every 4 weeks
	defer ticker.Stop()

	select {
	case <-ticker.C:
		// Operate if in enterprise mode
		if m.isEnterpriseMode() && licenseMode == "online" {
			// Fetch the public key periodically
			if err := m.fetchPublicKeyWithLock(); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Could not fetch public key for license file", err, nil)
				break
			}
		}
	}
}

func (m *Manager) fetchPublicKeyWithoutLock() error {
	// Check if offline licensing mode is used
	if licenseMode == "offline" {
		// Marshal the public key
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(licensePublicKey))
		if err != nil {
			return err
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
		return err
	}

	// Decode the response
	v := new(model.GraphqlFetchPublicKeyResponse)
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	// Check if valid response was received
	if v.Status != http.StatusOK {
		return fmt.Errorf("%s-%s", v.Message, v.Error)
	}

	// Marshal the public key
	publicKey := new(rsa.PublicKey)
	if err = json.Unmarshal([]byte(v.Result), publicKey); err != nil {
		return err
	}

	// Set the public key
	m.publicKey = publicKey
	return nil
}

func (m *Manager) ValidateLicense() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, err := m.decryptLicense(m.config.License); err != nil {
		m.ResetQuotas()
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to validate license key", err, nil)
	}

	return nil
}

func (m *Manager) RenewLicense(force bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.checkIfLeaderGateway() {
		return errors.New("only the leader can fetch the license")
	}

	// Throw error if licensing mode is set to offline
	if licenseMode == "offline" {
		return errors.New("cannot renew license in offline licensing mode")
	}

	return m.renewLicenseWithoutLock(force)
}

func (m *Manager) renewLicenseWithoutLock(force bool) error {
	// Marshal the request body
	data, _ := json.Marshal(map[string]interface{}{
		"params": model.RenewLicense{
			LicenseKey:       m.config.LicenseKey,
			LicenseValue:     m.config.LicenseValue,
			License:          m.config.License,
			CurrentSessionID: m.sessionID,
		},
		"timeout": 10,
	})
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), `Renewing admin license`, map[string]interface{}{"clusterId": m.config.LicenseKey, "clusterKey": m.config.LicenseValue, "sessionId": m.sessionID})
	// Fire the request
	res, err := http.Post("https://api.spaceuptech.com/v1/api/spacecloud/services/billing/renewLicense", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to contact server to fetch license file from server", err, nil)
	}
	defer func() { _ = res.Body.Close() }()

	// Decode the response
	data, _ = ioutil.ReadAll(res.Body)

	v := new(model.GraphqlFetchLicenseResponse)
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid status code received in response", errors.New(v.Error), nil)
	}

	// Check if response is valid
	if v.Status != http.StatusOK {
		m.licenseFetchErrorCount++
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to fetch license file. Retry count - %d", m.licenseFetchErrorCount), errors.New(v.Message), nil)
		if m.licenseFetchErrorCount > maxLicenseFetchErrorCount || force {
			helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), "Max retry limit hit.", nil)
			m.ResetQuotas()
			return fmt.Errorf("%s-%s", v.Message, v.Error)
		}
		return nil
	} else {
		m.licenseFetchErrorCount = 0
	}

	m.config.License = v.Result.License
	if err := m.setQuotas(v.Result.License); err != nil {
		return err
	}

	go func() { _ = m.syncMan.SetAdminConfig(context.TODO(), m.config) }()
	return nil
}

func (m *Manager) ResetQuotas() {
	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), "Resetting space cloud to run in open source model. You will have to re-register the cluster again.", nil)
	m.quotas.MaxProjects = 1
	m.quotas.MaxDatabases = 1
	m.quotas.IntegrationLevel = 0
	m.plan = "space-cloud-open--monthly"

	if licenseMode == "online" {
		m.config.LicenseKey = ""
		m.config.LicenseValue = ""
	}

	m.config.License = ""

	m.clusterName = ""

	go func() {
		if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to save admin config", err, nil)
		}
	}()
}

func (m *Manager) setQuotas(license string) error {
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to fetch public key", err, nil)
		}
	}
	licenseObj, err := m.decryptLicense(license)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to decrypt license key", err, nil)
	}

	if licenseMode == "offline" && m.sessionID != licenseObj.SessionID {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid license key provided. Make sure you use the license key for this cluster.", nil, nil)
		m.ResetQuotas()
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

func (m *Manager) isEnterpriseMode() bool {
	return m.isRegistered() && !strings.HasPrefix(m.plan, "space-cloud-open")
}

func (m *Manager) isRegistered() bool {
	return m.config.LicenseKey != "" && m.config.LicenseValue != "" && m.config.License != ""
}

func (m *Manager) decryptLicense(license string) (*model.License, error) {
	obj, err := m.parseLicenseToken(license)
	if err != nil {
		return nil, err
	}

	v := new(model.License)
	if err := mapstructure.Decode(obj, v); err != nil {
		return nil, err
	}
	return v, nil

}

func (m *Manager) parseLicenseToken(tokenString string) (map[string]interface{}, error) {
	licenseObj, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("invalid signing method type")
		}

		return m.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := licenseObj.Claims.(jwt.MapClaims); ok && licenseObj.Valid {
		if err := claims.Valid(); err != nil {
			return nil, err
		}
		return claims, nil
	}

	return nil, errors.New("unable to parse license token")
}

func (m *Manager) checkIfLeaderGateway() bool {
	return strings.HasSuffix(m.nodeID, "-0")
}
