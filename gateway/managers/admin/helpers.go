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

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
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
				if m.checkIfLeaderGateway() {
					// Fetch the public key periodically
					if err := m.RenewLicense(false); err != nil {
						_ = utils.LogError("Unable to renew license. Has your subscription expired?", "admin", "licenseRenewalRoutine", err)
						break
					}
					go func() {
						if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
							_ = utils.LogError("Unable to save admin config", "admin", "licenseRenewalRoutine", err)
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
		if m.isEnterpriseMode() {
			// Fetch the public key periodically
			if err := m.fetchPublicKeyWithLock(); err != nil {
				_ = utils.LogError("Could not fetch public key for license file", "admin", "fetch-license-routine", err)
				break
			}
		}
	}
}

func (m *Manager) fetchPublicKeyWithoutLock() error {
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
		return utils.LogError("Unable to validate license key", "admin", "ValidateLicense", err)
	}

	return nil
}

func (m *Manager) RenewLicense(force bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.checkIfLeaderGateway() {
		return errors.New("only the leader can fetch the license")
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
	utils.LogDebug(`Renewing admin license`, "admin", "renewLicenseWithoutLock", map[string]interface{}{"clusterId": m.config.LicenseKey, "clusterKey": m.config.LicenseValue, "sessionId": m.sessionID})
	// Fire the request
	res, err := http.Post("https://api.spaceuptech.com/v1/api/spacecloud/services/billing/renewLicense", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return utils.LogError("Unable to contact server to fetch license file from server", "admin", "renewLicenseWithoutLock", err)
	}
	defer func() { _ = res.Body.Close() }()

	// Decode the response
	data, _ = ioutil.ReadAll(res.Body)

	v := new(model.GraphqlFetchLicenseResponse)
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return utils.LogError("Invalid status code received in response", "admin", "renewLicenseWithoutLock", errors.New(v.Error))
	}

	// Check if response is valid
	if v.Status != http.StatusOK {
		m.licenseFetchErrorCount++
		_ = utils.LogError(fmt.Sprintf("Unable to fetch license file. Retry count - %d", m.licenseFetchErrorCount), "admin", "renewLicenseWithoutLock", errors.New(v.Message))
		if m.licenseFetchErrorCount > maxLicenseFetchErrorCount || force {
			utils.LogInfo("Max retry limit hit.", "admin", "renewLicenseWithoutLock")
			m.ResetQuotas()
			return fmt.Errorf("%s-%s", v.Message, v.Error)
		}
		return nil
	} else {
		m.licenseFetchErrorCount = 0
	}

	m.config.License = v.Result.License
	return m.setQuotas(v.Result.License)
}

func (m *Manager) ResetQuotas() {
	// TODO set sync man
	utils.LogInfo("Resetting space cloud to run in open source model. You will have to re-register the cluster again.", "admin", "resetQuotas")
	m.quotas.MaxProjects = 1
	m.quotas.MaxDatabases = 1
	m.quotas.IntegrationLevel = 0
	m.plan = "space-cloud-open--monthly"

	m.config.LicenseKey = ""
	m.config.LicenseValue = ""
	m.config.License = ""

	m.clusterName = ""

	go func() {
		if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
			_ = utils.LogError("Unable to save admin config", "admin", "resetQuotas", err)
		}
	}()
}

func (m *Manager) setQuotas(license string) error {
	if m.publicKey == nil {
		if err := m.fetchPublicKeyWithoutLock(); err != nil {
			return utils.LogError("Unable to fetch public key", "admin", "setQuotas", err)
		}
	}
	licenseObj, err := m.decryptLicense(license)
	if err != nil {
		return utils.LogError("Unable to decrypt license key", "admin", "setQuotas", err)
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
	return m.config.LicenseKey != "" && m.config.LicenseValue != ""
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
