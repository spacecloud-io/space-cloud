package ingress

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/filter"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

func deleteIngressGlobalConfig(project string) error {

	// Delete the ingress global config from the server
	url := fmt.Sprintf("/v1/config/projects/%s/routing/ingress/global", project)

	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, payload); err != nil {
		return err
	}

	return nil
}

func deleteIngressRoute(project, prefix string) error {

	objs, err := GetIngressRoutes(project, "ingress-route", map[string]string{})
	if err != nil {
		return err
	}

	doesRouteExist := false
	routeIDs := []string{}
	for _, spec := range objs {
		routeIDs = append(routeIDs, spec.Meta["id"])
	}

	prefix, err = filter.DeleteOptions(prefix, routeIDs, doesRouteExist)
	if err != nil {
		return err
	}

	// Delete the ingress route from the server
	url := fmt.Sprintf("/v1/config/projects/%s/routing/ingress/%s", project, prefix)

	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodDelete, url, map[string]string{}, payload); err != nil {
		return err
	}

	return nil
}
