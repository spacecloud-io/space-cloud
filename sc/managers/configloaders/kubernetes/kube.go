package kube

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/managers/configloaders/common"
	"github.com/spacecloud-io/space-cloud/model"
	"go.uber.org/zap"
)

var connectorPool = caddy.NewUsagePool()

func init() {
	caddy.RegisterModule(Loader{})
}

// Loader loads the sc configuration from config maps
type Loader struct {
	logger *zap.Logger

	connector *Connector
}

// CaddyModule returns the Caddy module information.
func (Loader) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.config_loaders.kube",
		New: func() caddy.Module { return new(Loader) },
	}
}

// Provision sets up the kube loader module.
func (k *Loader) Provision(ctx caddy.Context) error {
	k.logger = ctx.Logger(k)

	poolKey := "kube-store"
	val, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return New()
	})
	if err != nil {
		k.logger.Error("Unable to load kube connector from connectorPool", zap.Error(err))
		return err
	}

	k.connector = val.(*Connector)

	return nil
}

// Cleanup handled cleanup when module is reloaded
func (k *Loader) Cleanup() error {
	poolKey := "kube-store"
	_, err := connectorPool.Delete(poolKey)
	if err != nil {
		k.logger.Error("Unable to cleanup kube connector from connector pool", zap.Error(err))
		return err
	}

	return nil
}

// LoadConfig returns the final caddy config from the store.
func (k *Loader) LoadConfig(ctx caddy.Context) ([]byte, error) {
	kubeConnector := k.connector

	// Load the new caddy config
	kubeConnector.Lock.Lock()
	defer kubeConnector.Lock.Unlock()

	// TODO: pass SC config to this function
	config, err := common.PrepareConfig(&model.SCConfig{})
	if err != nil {
		k.logger.Error("Unable to prepare caddy config", zap.Error(err))
		return nil, err
	}

	return json.Marshal(config)
}

// Interface guards
var (
	_ caddy.Provisioner  = (*Loader)(nil)
	_ caddy.ConfigLoader = (*Loader)(nil)
	_ caddy.CleanerUpper = (*Loader)(nil)
)
