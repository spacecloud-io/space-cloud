package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/modules/auth/types"
	"github.com/spacecloud-io/space-cloud/utils"
)

// KratosAuthVerifyHandler is responsible to authenticate the incoming request
// using Kratos
type KratosAuthVerifyHandler struct {
	logger *zap.Logger
	//authApp *App
}

// CaddyModule returns the Caddy module information.
func (KratosAuthVerifyHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.sc_kratos_auth_verify_handler",
		New: func() caddy.Module { return new(KratosAuthVerifyHandler) },
	}
}

// Provision sets up the auth verify module.
func (h *KratosAuthVerifyHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)

	return nil
}

// ServeHTTP handles the http request
func (h *KratosAuthVerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Prepare authentication object
	result := types.AuthResult{}
	// Check if token is present in the header

	var foo map[string]interface{}
	client := &http.Client{}

	kratosEndpoint := viper.GetString("kratos.endpoint")
	// Create a new GET request to the /sessions/whoami endpoint
	req, err := http.NewRequest("GET", kratosEndpoint+"/sessions/whoami", nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return err
	}
	if cookie := r.Header.Get("Cookie"); len(cookie) != 0 {
		req.Header.Set("Cookie", cookie)
	}
	if cookie := r.Header.Get("Cookie"); len(cookie) != 0 {
		req.Header.Set("Cookie", cookie)
	}
	if token, p := getTokenFromHeader(r); p {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	//token := "ory_st_94jGcnytycUNCtKIVchOfiyWrs8wlJYY"
	//req.Header.Set("Authorization", "Bearer "+token)
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	err = json.NewDecoder(resp.Body).Decode(&foo)
	//body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}
	// Print the response body
	//fmt.Println(string(body))
	j, _ := json.MarshalIndent(foo, "", " ")
	fmt.Println(string(j))

	// Add the result in the context object
	r = utils.StoreAuthenticationResult(r, &result)
	return next.ServeHTTP(w, r)
}

// Interface guard
var _ caddy.Provisioner = (*KratosAuthVerifyHandler)(nil)
var _ caddyhttp.MiddlewareHandler = (*KratosAuthVerifyHandler)(nil)
