package admin

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (m *Manager) setPublicCloseChannel(ch chan struct{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.closeFetchPublicRSAKey = ch
}

func (m *Manager) fetchPublicKeyWithLock() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.fetchPublicKeyWithoutLock()
}

func (m *Manager) fetchPublicKeyRoutine() {
	// Create a new ticker
	ticker := time.NewTicker(4 * 7 * 24 * time.Hour) // fetch public once every 4 weeks
	defer ticker.Stop()

	// Make a new closer channel
	closeCh := make(chan struct{})
	defer close(closeCh)
	m.setPublicCloseChannel(closeCh)

	select {
	case <-ticker.C:
		// Fetch the public key periodically
		err := m.fetchPublicKeyWithLock()
		if err != nil {
			logrus.Errorf("Could not fetch public key from spaceuptech server - %s", err.Error())
		}

	case <-closeCh:
		// Close the routine on receiving the command
		return
	}
}

func (m *Manager) fetchPublicKeyWithoutLock() error {
	// Fire the http request
	res, err := http.Get(fmt.Sprintf("http://localhost:4100/v1/public-key"))
	if err != nil {
		return err
	}

	// Decode the response
	data := new(model.Response)
	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		return err
	}

	// Check if valid response was received
	if res.StatusCode != http.StatusOK {
		return errors.New(data.Error)
	}

	// Marshal the public key
	publicKey := new(rsa.PublicKey)
	if err = json.Unmarshal([]byte(data.Result.(string)), publicKey); err != nil {
		return err
	}

	// Set the public key
	m.publicKey = publicKey
	return nil
}

func (m *Manager) fetchQuotas() error {
	// Marshal the request pody
	data, _ := json.Marshal(map[string]string{"clusterId": m.config.ClusterID, "clusterKey": m.config.ClusterKey})

	// Fire the request
	res, err := http.Post("http://localhost:4100/v1/quotas", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Decode the response
	v := new(model.UsageQuotasResult)
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		logrus.Println("error", err)
		return err
	}

	// Check if response is valid
	if res.StatusCode != http.StatusOK {
		return errors.New(v.Error)
	}

	// Set the quotas and version
	m.quotas = v.Result
	return nil
}

func (m *Manager) isEnterpriseMode() bool {
	return m.config.ClusterID != "" && m.config.ClusterKey != ""
}
