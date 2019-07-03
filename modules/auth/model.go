package auth

import "errors"

// TokenClaims holds the JWT token claims
type TokenClaims map[string]interface{}

// GetRole returns the role present in the token claims
func (c TokenClaims) GetRole() (string, error) {
	roleTemp, p := c["role"]
	if !p {
		return "", errors.New("role is not present in the token claims")
	}

	role, ok := roleTemp.(string)
	if !ok {
		return "", errors.New("role is not of the correct type")
	}

	return role, nil
}
