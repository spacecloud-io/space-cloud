package utils

import (
	"context"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Auth holds the required fields for jwt package
type Auth struct {
	lock   sync.RWMutex
	secret string
}

func New(secret string) *Auth {
	return &Auth{
		secret: secret,
	}
}

// CreateToken create a token with primary secret
func (s *Auth) CreateToken(ctx context.Context, tokenClaims map[string]interface{}) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	claims := jwt.MapClaims{}
	for k, v := range tokenClaims {
		claims[k] = v
	}
	// Add expiry of one week
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}
