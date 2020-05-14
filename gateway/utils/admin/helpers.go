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

	"github.com/spaceuptech/space-cloud/gateway/config"
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
						_ = utils.LogError("Unable to renew license. Has your subscription expired?", err)
						break
					}
					go func() {
						if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
							_ = utils.LogError("Unable to save admin config", err)
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

// func (m *Manager) fetchPublicKeyRoutine() {
// 	// Create a new ticker
// 	ticker := time.NewTicker(4 * 7 * 24 * time.Hour) // fetch public once every 4 weeks
// 	defer ticker.Stop()
//
// 	select {
// 	case <-ticker.C:
// 		// Operate if in enterprise mode
// 		if m.isEnterpriseMode() {
// 			// Fetch the public key periodically
// 			if err := m.fetchPublicKeyWithLock(); err != nil {
// 				_ = utils.LogError("Could not fetch public key for license file", err)
// 				break
// 			}
//
// 			if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
// 				_ = utils.LogError("Unable to save admin config", err)
// 			}
// 		}
// 	}
// }

func (m *Manager) fetchPublicKeyWithoutLock() error {
	// Fire the http request
	body := map[string]interface{}{
		"timeout": 10,
	}
	data, _ := json.Marshal(body)
	res, err := http.Post(fmt.Sprintf("http://35.188.208.249/v1/api/spacecloud/services/backend/fetch_public_key"), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Decode the response
	v := new(model.GraphqlFetchPublicKeyResponse)
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	// Check if valid response was received
	if v.Result.Status != http.StatusOK {
		return fmt.Errorf("%s-%s", v.Result.Message, v.Result.Error)
	}

	// Marshal the public key
	publicKey := new(rsa.PublicKey)
	if err = json.Unmarshal([]byte(v.Result.Result), publicKey); err != nil {
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
		m.resetQuotas()
		return utils.LogError("Unable to validate license key", err)
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
			ClusterID:        m.config.ClusterID,
			ClusterKey:       m.config.ClusterKey,
			License:          m.config.License,
			CurrentSessionID: m.sessionID,
		},
		"timeout": 10,
	})
	utils.LogDebug(`Renewing admin license`, map[string]interface{}{"clusterId": m.config.ClusterID, "clusterKey": m.config.ClusterKey, "sessionId": m.sessionID})
	// Fire the request
	res, err := http.Post("http://35.188.208.249/v1/api/spacecloud/services/backend/fetch_license", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return utils.LogError("Unable to contact server to fetch license file from server", err)
	}
	defer func() { _ = res.Body.Close() }()

	// Decode the response
	data, _ = ioutil.ReadAll(res.Body)

	v := new(model.GraphqlFetchLicenseResponse)
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return utils.LogError("Invalid status code received in response", errors.New(v.Error))
	}

	// Check if response is valid
	if v.Result.Status != http.StatusOK {
		m.licenseFetchErrorCount++
		_ = utils.LogError(fmt.Sprintf("Unable to fetch license file. Retry count - %d", m.licenseFetchErrorCount), errors.New(v.Result.Message))
		if m.licenseFetchErrorCount > maxLicenseFetchErrorCount || force {
			utils.LogInfo("Max retry limit hit.")
			m.resetQuotas()
			return fmt.Errorf("%s-%s", v.Result.Message, v.Result.Error)
		}
		return nil
	} else {
		m.licenseFetchErrorCount = 0
	}

	m.config.License = v.Result.Result.License
	return m.setQuotas(v.Result.Result.License)
}

func (m *Manager) resetQuotas() {
	// TODO set sync man
	utils.LogInfo("Resetting space cloud to run in open source model. You will have to re-register the cluster again.")
	m.quotas.MaxProjects = 1
	m.quotas.MaxDatabases = 1
	m.plan = "space-cloud-open--monthly"
	m.config = new(config.Admin)

	go func() {
		if err := m.syncMan.SetAdminConfig(context.Background(), m.config); err != nil {
			_ = utils.LogError("Unable to save admin config", err)
		}
	}()
}

func (m *Manager) setQuotas(license string) error {
	if err := m.fetchPublicKeyWithoutLock(); err != nil {
		return utils.LogError("Unable to fetch public key", err)
	}
	licenseObj, err := m.decryptLicense(license)
	if err != nil {
		return utils.LogError("Unable to decrypt license key", err)
	}

	// set quotas
	m.quotas.MaxProjects = licenseObj.Quotas.MaxProjects
	m.quotas.MaxDatabases = licenseObj.Quotas.MaxDatabases
	m.plan = licenseObj.Plan
	return nil
}
func (m *Manager) isEnterpriseMode() bool {
	return m.isRegistered() && !strings.HasPrefix(m.plan, "space-cloud-open")
}

func (m *Manager) isRegistered() bool {
	return m.config.ClusterID != "" && m.config.ClusterKey != ""
}

func (m *Manager) decryptLicense(license string) (*model.License, error) {
	licenseObj, err := jwt.Parse(license, func(token *jwt.Token) (interface{}, error) {
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

		v := new(model.License)
		if err := mapstructure.Decode(claims, v); err != nil {
			return nil, err
		}
		return v, nil
	}

	return nil, errors.New("unable to parse license")
}

func (m *Manager) checkIfLeaderGateway() bool {
	return strings.HasSuffix(m.nodeID, "-0")
}
