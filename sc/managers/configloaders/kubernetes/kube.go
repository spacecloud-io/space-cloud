package kubernetes

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/configloaders/common"
	"github.com/spacecloud-io/space-cloud/managers/configloaders/kubernetes/kubeconnector"
)

var connectorPool = caddy.NewUsagePool()

func init() {
	caddy.RegisterModule(KubeLoader{})
}

// KubeLoader loads the sc configuration from config maps
type KubeLoader struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (KubeLoader) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.config_loaders.kube",
		New: func() caddy.Module { return new(KubeLoader) },
	}
}

// Provision sets up the kube loader module.
func (k *KubeLoader) Provision(ctx caddy.Context) error {
	k.logger = ctx.Logger(k)

	poolKey := "kube-store"
	_, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return kubeconnector.New()
	})
	if err != nil {
		k.logger.Error("Unable to load kube connector from connectorPool", zap.Error(err))
		return err
	}

	return nil
}

// Cleanup handled cleanup when module is reloaded
func (k *KubeLoader) Cleanup() error {
	poolKey := "kube-store"
	_, err := connectorPool.Delete(poolKey)
	if err != nil {
		k.logger.Error("Unable to cleanup kube connector from connector pool", zap.Error(err))
		return err
	}

	return nil
}

// LoadConfig returns the final caddy config from the store.
func (k *KubeLoader) LoadConfig(ctx caddy.Context) ([]byte, error) {

	poolKey := "kube-store"
	val, _, err := connectorPool.LoadOrNew(poolKey, func() (caddy.Destructor, error) {
		return kubeconnector.New()
	})
	if err != nil {
		k.logger.Error("Unable to load kube connector", zap.Error(err))
		return nil, err
	}

	kubeConnector := val.(*kubeconnector.Connector)

	// Load the new caddy config
	kubeConnector.Lock.Lock()
	defer kubeConnector.Lock.Unlock()
	config, err := common.PrepareConfig(kubeConnector.ProjectsConfig)
	if err != nil {
		k.logger.Error("Unable to prepare caddy config", zap.Error(err))
		return nil, err
	}

	return json.Marshal(config)
}

// Interface guards
var (
	_ caddy.Provisioner  = (*KubeLoader)(nil)
	_ caddy.ConfigLoader = (*KubeLoader)(nil)
	_ caddy.CleanerUpper = (*KubeLoader)(nil)
)
