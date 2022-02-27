package common

import (
	"encoding/json"
	"strconv"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/spf13/viper"
)

func prepareHTTPHanndlerApp() []byte {
	port := viper.GetInt("port")
	sshCert := viper.GetString("ssl-cert")
	sshKey := viper.GetString("ssl-key")

	httpsPort := 0
	listen := []string{":" + strconv.Itoa(port)}
	if sshCert != "none" && sshKey != "none" {
		httpsPort = port + 4
		listen = []string{":" + strconv.Itoa(httpsPort)}
		port = 0
	}

	httpConfig := caddyhttp.App{
		HTTPPort:  port,
		HTTPSPort: httpsPort,
		Servers: map[string]*caddyhttp.Server{
			"default": {
				Listen: listen,
				Routes: caddyhttp.RouteList{
					caddyhttp.Route{
						HandlersRaw: []json.RawMessage{getHandler()},
					},
				},
			},
		},
	}

	data, _ := json.Marshal(httpConfig)
	return data
}

func getHandler() []byte {
	handler := make(map[string]string)

	handler["handler"] = "sc_handler"

	data, _ := json.Marshal(handler)
	return data
}
