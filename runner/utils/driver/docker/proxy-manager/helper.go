package manager

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

func getConfigKey(projectID, serviceID string) string {
	return projectID + "-" + serviceID
}

func getProjectAndServiceIDFromKey(key string) (projectID, serviceID string) {
	array := strings.Split(key, "-")
	return array[0], array[1]
}
func getHost(req *http.Request) string {
	return strings.Split(req.Host, ":")[0]
}

func getProjectAndServiceID(host string) (projectID, serviceID string) {
	array := strings.Split(host, ".")
	return array[1], array[0]
}

func getServiceAndProject(req *http.Request) (string, string) {
	host := getHost(req)
	return getProjectAndServiceID(host)
}

func (m *Manager) getRoute(projectID, serviceID string, port int32) (*model.Route, error) {

	// Check if the service id exists
	routes, p := m.serviceRoutes[getConfigKey(projectID, serviceID)]
	if !p {
		return nil, fmt.Errorf("no routes found for service (%s) in project (%s)", serviceID, projectID)
	}

	// Select the correct route based on the port
	for _, route := range routes {
		if route.Source.Port == port {
			return route, nil
		}
	}
	return nil, fmt.Errorf("no routes found for port (%d) for service (%s) in project (%s)", port, serviceID, projectID)
}

func setRequest(request *http.Request, route *model.Route, projectID, serviceID string) error {
	// http: Request.RequestURI can't be set in client requests.
	// http://golang.org/src/pkg/net/http/client.go
	request.RequestURI = ""

	// Change the request with the destination host, port and url
	target, err := route.SelectTarget(-1) // pass a -ve weight to randomly generate
	if err != nil {
		return err
	}

	host := target.Host
	if target.Type == model.RouteTargetVersion {
		host = utils.GetInternalServiceDomain(projectID, serviceID, target.Version)
	}

	request.Host = host
	request.URL.Host = fmt.Sprintf("%s:%d", host, target.Port)

	// Set the url scheme to http
	if target.Scheme == "" {
		target.Scheme = "http"
	}
	request.URL.Scheme = target.Scheme
	return nil
}

func (m *Manager) adjustProxyServers() {

	// Calculate the ports required
	newPorts := make(map[int32]struct{})
	for _, routes := range m.serviceRoutes {
		// Make a map of ports requested!
		for r := range routes {
			for _, p := range routes[r].Targets {
				newPorts[p.Port] = struct{}{}
			}
		}
	}

	// Check for the ports to be closed
	for port, server := range m.servers {
		if _, p := newPorts[port]; !p {
			_ = server.Close()
			delete(m.servers, port)
		}
	}

	// Check for the ports to be started
	for port := range newPorts {
		if _, p := m.servers[port]; !p {
			m.startServer(port)
		}
	}
}

func (m *Manager) startServer(port int32) {
	obj := &http.Server{Addr: ":" + strconv.Itoa(int(port)), Handler: m.routes(port)}
	go func() { _ = obj.ListenAndServe() }()
	m.servers[port] = obj
}
