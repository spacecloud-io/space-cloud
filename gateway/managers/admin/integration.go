package admin

// GetIntegrationToken returns the admin token required by an intergation
func (m *Manager) GetIntegrationToken(id string) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.createToken(map[string]interface{}{"id": id, "role": "integration"})
}