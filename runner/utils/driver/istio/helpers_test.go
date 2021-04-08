package istio

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/gogo/protobuf/types"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func Test_prepareVirtualServiceTCPRoutes(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		serviceID string
		services  map[string]model.AutoScaleConfig
		routes    model.Routes
		proxyPort uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []*networkingv1alpha3.TCPRoute
		wantErr bool
	}{
		{
			name: "Multiple routes & multiple targets for TCP protocol",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.TCP,
							Port:     8080,
						},
						Targets: []model.RouteTarget{
							{
								Port:    8080,
								Weight:  25,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:    8080,
								Weight:  25,
								Version: "v3",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:   8080,
								Weight: 25,
								Type:   model.RouteTargetExternal,
								Host:   "httpbin.test.svc.cluster.local",
							},
							{
								Port:   8080,
								Weight: 25,
								Type:   model.RouteTargetExternal,
								Host:   "httpbin.test.svc.cluster.local",
							},
						},
					},
				},
				proxyPort: 4050,
			},
			want: []*networkingv1alpha3.TCPRoute{
				{
					Match: []*networkingv1alpha3.L4MatchAttributes{{Port: uint32(8080)}},
					Route: []*networkingv1alpha3.RouteDestination{
						{
							Destination: &networkingv1alpha3.Destination{
								Host: getInternalServiceDomain("myProject", "greeter", "v2"),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
						{
							Destination: &networkingv1alpha3.Destination{
								Host: getInternalServiceDomain("myProject", "greeter", "v3"),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
						{
							Destination: &networkingv1alpha3.Destination{
								Host: "httpbin.test.svc.cluster.local",
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
						{
							Destination: &networkingv1alpha3.Destination{
								Host: "httpbin.test.svc.cluster.local",
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Don't create TCP routes when protocol is HTTP",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.HTTP,
							Port:     8080,
						},
						Targets: []model.RouteTarget{
							{
								Port:    8080,
								Weight:  50,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:    8080,
								Weight:  50,
								Version: "v3",
								Type:    model.RouteTargetVersion,
							},
						},
					},
				},
				proxyPort: 4050,
			},
			want:    nil,
			wantErr: false,
		},
		// Check errors
		{
			name: "Source port not provided",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.TCP,
							Port:     0,
						},
						Targets: []model.RouteTarget{
							{
								Port:    8080,
								Weight:  50,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:    8080,
								Weight:  50,
								Version: "v3",
								Type:    model.RouteTargetVersion,
							},
						},
					},
				},
				proxyPort: 4050,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Route target not provided",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.TCP,
							Port:     0,
						},
						Targets: nil,
					},
				},
				proxyPort: 4050,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unknown service version provided in route target",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.TCP,
							Port:     0,
						},
						Targets: []model.RouteTarget{
							{
								Port:    8080,
								Weight:  50,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:    8080,
								Weight:  50,
								Version: "v4",
								Type:    model.RouteTargetVersion,
							},
						},
					},
				},
				proxyPort: 4050,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unknown route type provided",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.TCP,
							Port:     0,
						},
						Targets: []model.RouteTarget{
							{
								Port:    8080,
								Weight:  50,
								Version: "v2",
								Type:    "interior",
							},
							{
								Port:    8080,
								Weight:  50,
								Version: "v4",
								Type:    model.RouteTargetVersion,
							},
						},
					},
				},
				proxyPort: 4050,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareVirtualServiceTCPRoutes(tt.args.ctx, tt.args.projectID, tt.args.serviceID, tt.args.services, tt.args.routes)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareVirtualServiceTCPRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if arr := deep.Equal(got, tt.want); len(arr) > 0 {
				a, _ := json.MarshalIndent(arr, "", " ")
				t.Errorf("prepareVirtualServiceTCPRoutes() diff = %v", string(a))
			}
		})
	}
}

func Test_prepareVirtualServiceHTTPRoutes(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		serviceID string
		services  map[string]model.AutoScaleConfig
		routes    model.Routes
		proxyPort uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []*networkingv1alpha3.HTTPRoute
		wantErr bool
	}{
		{
			name: "Multiple routes & multiple targets for TCP protocol",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
				serviceID: "greeter",
				services: map[string]model.AutoScaleConfig{
					"v2": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
					"v3": {
						PollingInterval:  15,
						CoolDownInterval: 120,
						MinReplicas:      1,
						MaxReplicas:      100,
					},
				},
				routes: []*model.Route{
					{
						ID: "greeter",
						Source: model.RouteSource{
							Protocol: model.HTTP,
							Port:     8080,
						},
						Targets: []model.RouteTarget{
							{
								Port:    8080,
								Weight:  25,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:    8080,
								Weight:  25,
								Version: "v3",
								Type:    model.RouteTargetVersion,
							},
							{
								Port:   8080,
								Weight: 25,
								Type:   model.RouteTargetExternal,
								Host:   "httpbin.test.svc.cluster.local",
							},
							{
								Port:   8080,
								Weight: 25,
								Type:   model.RouteTargetExternal,
								Host:   "httpbin.test.svc.cluster.local",
							},
						},
						Matchers: []*model.Matcher{
							{
								URL: &model.HTTPMatcher{
									Value:      "/v2/",
									Type:       model.RouteHTTPMatchTypeExact,
									IgnoreCase: true,
								},
								Headers: []*model.HTTPMatcher{
									{
										Key:   "version",
										Value: "v2",
										Type:  model.RouteHTTPMatchTypeExact,
									},
								},
							},
							{
								URL: &model.HTTPMatcher{
									Value: "/v2/",
									Type:  model.RouteHTTPMatchTypePrefix,
								},
								Headers: []*model.HTTPMatcher{
									{
										Key:   "version",
										Value: "v2",
										Type:  model.RouteHTTPMatchTypePrefix,
									},
								},
							},
							{
								URL: &model.HTTPMatcher{
									Value: "/v2/",
									Type:  model.RouteHTTPMatchTypeRegex,
								},
								Headers: []*model.HTTPMatcher{
									{
										Key:   "version",
										Value: "v2",
										Type:  model.RouteHTTPMatchTypeRegex,
									},
								},
							},
							{
								Headers: []*model.HTTPMatcher{
									{
										Key:  "version",
										Type: model.RouteHTTPMatchTypeCheckPresence,
									},
								},
							},
						},
					},
				},
				proxyPort: 4050,
			},
			want: []*networkingv1alpha3.HTTPRoute{
				{
					Name: "http-8080",
					Match: []*networkingv1alpha3.HTTPMatchRequest{
						{
							Uri:           &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Exact{Exact: "/v2/"}},
							Headers:       map[string]*networkingv1alpha3.StringMatch{"version": {MatchType: &networkingv1alpha3.StringMatch_Exact{Exact: "v2"}}},
							IgnoreUriCase: true,
							Port:          uint32(8080),
							Gateways:      []string{"mesh"},
						},
						{
							Uri:      &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Prefix{Prefix: "/v2/"}},
							Headers:  map[string]*networkingv1alpha3.StringMatch{"version": {MatchType: &networkingv1alpha3.StringMatch_Prefix{Prefix: "v2"}}},
							Port:     uint32(8080),
							Gateways: []string{"mesh"},
						},
						{
							Uri:      &networkingv1alpha3.StringMatch{MatchType: &networkingv1alpha3.StringMatch_Regex{Regex: "/v2/"}},
							Headers:  map[string]*networkingv1alpha3.StringMatch{"version": {MatchType: &networkingv1alpha3.StringMatch_Regex{Regex: "v2"}}},
							Port:     uint32(8080),
							Gateways: []string{"mesh"},
						},
						{
							Headers:  map[string]*networkingv1alpha3.StringMatch{"version": {MatchType: nil}},
							Port:     uint32(8080),
							Gateways: []string{"mesh"},
						},
					},
					Retries: &networkingv1alpha3.HTTPRetry{Attempts: 3, PerTryTimeout: &types.Duration{Seconds: 180}},
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						{
							Headers: &networkingv1alpha3.Headers{
								Request: &networkingv1alpha3.Headers_HeaderOperations{
									Set: map[string]string{
										"x-og-project": "myProject",
										"x-og-service": "greeter",
										"x-og-host":    "greeter-v2-internal.myProject.svc.cluster.local",
										"x-og-port":    "8080",
										"x-og-version": "v2",
									},
								},
							},
							Destination: &networkingv1alpha3.Destination{
								Host: getInternalServiceDomain("myProject", "greeter", "v2"),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
						{
							Headers: &networkingv1alpha3.Headers{
								Request: &networkingv1alpha3.Headers_HeaderOperations{
									Set: map[string]string{
										"x-og-project": "myProject",
										"x-og-service": "greeter",
										"x-og-host":    "greeter-v3-internal.myProject.svc.cluster.local",
										"x-og-port":    "8080",
										"x-og-version": "v3",
									},
								},
							},
							Destination: &networkingv1alpha3.Destination{
								Host: getInternalServiceDomain("myProject", "greeter", "v3"),
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
						{
							Destination: &networkingv1alpha3.Destination{
								Host: "httpbin.test.svc.cluster.local",
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
						{
							Destination: &networkingv1alpha3.Destination{
								Host: "httpbin.test.svc.cluster.local",
								Port: &networkingv1alpha3.PortSelector{Number: uint32(8080)},
							},
							Weight: 25,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareVirtualServiceHTTPRoutes(tt.args.ctx, tt.args.projectID, tt.args.serviceID, tt.args.services, tt.args.routes, tt.args.proxyPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareVirtualServiceHTTPRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if arr := deep.Equal(got, tt.want); len(arr) > 0 {
				a, _ := json.MarshalIndent(arr, "", " ")
				t.Errorf("prepareVirtualServiceTCPRoutes() diff = %v", string(a))
			}
		})
	}
}
