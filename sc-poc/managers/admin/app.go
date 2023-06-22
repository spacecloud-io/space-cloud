package admin

import (
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(App{})
}

// App describes the source manager app
type App struct {
	// Internal stuff
	logger   *zap.Logger
	username string
	password string
	secret   string
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "sc_admin",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the source manager.
func (a *App) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	a.username = viper.GetString("admin.username")
	a.password = viper.GetString("admin.password")
	a.secret = viper.GetString("admin.secret")
	return nil
}

func (a *App) SignSCToken() (string, error) {
	jwtClaims := make(jwt.MapClaims)
	jwtClaims["sub"] = a.username
	jwtClaims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *App) VerifyCredentials(creds map[string]string) error {
	if creds["username"] != a.username || creds["password"] != a.password {
		return fmt.Errorf("invalid credentials provided")
	}
	return nil
}

func (a *App) VerifySCToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("invalid token found in the Authorization Header")
	}

	return a.SignSCToken()
}

// Start begins the source manager operations
func (a *App) Start() error {
	return nil
}

// Stop ends the source manager operations
func (a *App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
