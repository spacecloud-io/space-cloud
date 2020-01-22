package auth

import "github.com/dgrijalva/jwt-go"

func (m *Module) GenerateHS256Token(serviceId, projectId, version string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      serviceId,
		"project": projectId,
		"version": version,
	})
	return token.SignedString([]byte(m.config.Secret))
}
