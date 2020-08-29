package routing

import (
	"context"
	"fmt"
	"text/template"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	tmpl2 "github.com/spaceuptech/space-cloud/gateway/utils/tmpl"
)

func (r *Routing) createGoTemplate(kind, project, id, tmpl string) error {
	key := getGoTemplateKey(kind, project, id)

	// Create a new template object
	t := template.New(key)
	t = t.Funcs(tmpl2.CreateGoFuncMaps(nil))
	val, err := t.Parse(tmpl)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid golang template provided", err, nil)
	}

	r.goTemplates[key] = val
	return nil
}

func getGoTemplateKey(kind, project, id string) string {
	return fmt.Sprintf("%s---%s---%s", project, id, kind)
}

func (r *Routing) adjustBody(ctx context.Context, kind, project, token string, route *config.Route, auth, params interface{}) (interface{}, error) {
	var req interface{}
	var err error

	switch route.Modify.Tmpl {
	case config.EndpointTemplatingEngineGo:
		if tmpl, p := r.goTemplates[getGoTemplateKey(kind, project, route.ID)]; p {
			req, err = tmpl2.GoTemplate(module, "go-template", tmpl, route.Modify.OpFormat, token, auth, params)
			if err != nil {
				return nil, err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", route.Modify.Tmpl), map[string]interface{}{"project": project, "kind": kind, "route": route.ID})
		return params, nil
	}

	if req == nil {
		return params, nil
	}
	return req, nil
}
