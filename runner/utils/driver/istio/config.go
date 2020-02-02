package istio

// Config describes the configuration used by the istio driver
type Config struct {
	IsInsideCluster bool
	KubeConfigPath  string
	ProxyPort       uint32
}

// GenerateInClusterConfig returns a in-cluster config
func GenerateInClusterConfig() *Config {
	return &Config{IsInsideCluster: true}
}

// GenerateOutsideClusterConfig returns an out-of-cluster config
func GenerateOutsideClusterConfig(kubeConfigPath string) *Config {
	return &Config{IsInsideCluster: false, KubeConfigPath: kubeConfigPath}
}

// SetProxyPort sets the port of the proxy runner
func (c *Config) SetProxyPort(port uint32) {
	c.ProxyPort = port
}
