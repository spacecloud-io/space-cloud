package istio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-test/deep"
	"github.com/gogo/protobuf/types"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubefake "k8s.io/client-go/kubernetes/fake"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func TestIstio_GetServiceRoutes(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name                      string
		args                      args
		virtualServiceToBeCreated *v1alpha3.VirtualService
		want                      map[string]model.Routes
		wantErr                   bool
	}{
		{
			name: "Get TCP Routes with internal & external targets",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
			},
			virtualServiceToBeCreated: &v1alpha3.VirtualService{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: getVirtualServiceName("greeter"),
					Labels: map[string]string{
						"app":                          "greeter",
						"app.kubernetes.io/name":       "greeter",
						"app.kubernetes.io/managed-by": "space-cloud",
						"space-cloud.io/version":       model.Version,
					},
				},
				Spec: networkingv1alpha3.VirtualService{
					Hosts: []string{getServiceDomainName("myProject", "greeter")},
					Tcp: []*networkingv1alpha3.TCPRoute{
						{
							Match: []*networkingv1alpha3.L4MatchAttributes{{Port: uint32(8080)}},
							Route: []*networkingv1alpha3.RouteDestination{
								{
									Destination: &networkingv1alpha3.Destination{
										Host: getInternalServiceDomain("myProject", "greeter", "v2"),
										Port: &networkingv1alpha3.PortSelector{Number: 9090},
									},
									Weight: 50,
								},
								{
									Destination: &networkingv1alpha3.Destination{
										Host: "httpbin.myProject.svc.cluster.local",
										Port: &networkingv1alpha3.PortSelector{Number: 1010},
									},
									Weight: 50,
								},
							},
						},
					},
				},
			},
			want: map[string]model.Routes{
				"greeter": []*model.Route{
					{
						ID:             "greeter",
						RequestRetries: 0,
						RequestTimeout: 0,
						Source: model.RouteSource{
							Protocol: model.TCP,
							Port:     8080,
						},
						Targets: []model.RouteTarget{
							{
								Host:    "",
								Port:    9090,
								Weight:  50,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Host:   "httpbin.myProject.svc.cluster.local",
								Port:   1010,
								Weight: 50,
								Type:   model.RouteTargetExternal,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Get HTTP Routes with internal & external targets",
			args: args{
				ctx:       context.Background(),
				projectID: "myProject",
			},
			virtualServiceToBeCreated: &v1alpha3.VirtualService{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: getVirtualServiceName("greeter"),
					Labels: map[string]string{
						"app":                          "greeter",
						"app.kubernetes.io/name":       "greeter",
						"app.kubernetes.io/managed-by": "space-cloud",
						"space-cloud.io/version":       model.Version,
					},
				},
				Spec: networkingv1alpha3.VirtualService{
					Hosts: []string{getServiceDomainName("myProject", "greeter")},
					Http: []*networkingv1alpha3.HTTPRoute{
						{
							Name:    fmt.Sprintf("http-%d", 8080),
							Match:   []*networkingv1alpha3.HTTPMatchRequest{{Port: uint32(8080), Gateways: []string{"mesh"}}},
							Retries: &networkingv1alpha3.HTTPRetry{Attempts: 5, PerTryTimeout: &types.Duration{Seconds: model.DefaultRequestTimeout}},
							Route: []*networkingv1alpha3.HTTPRouteDestination{
								{
									Headers: &networkingv1alpha3.Headers{
										Request: &networkingv1alpha3.Headers_HeaderOperations{
											Set: map[string]string{
												"x-og-project": "myProject",
												"x-og-service": "greeter",
												"x-og-host":    getInternalServiceDomain("myProject", "greeter", "v2"),
												"x-og-port":    strconv.Itoa(int(9090)),
												"x-og-version": "v2",
											},
										},
									},
									Destination: &networkingv1alpha3.Destination{
										Host: getInternalServiceDomain("myProject", "greeter", "v2"),
										Port: &networkingv1alpha3.PortSelector{Number: 9090},
									},
									Weight: 50,
								},
								{
									Destination: &networkingv1alpha3.Destination{
										Host: "httpbin.myProject.svc.cluster.local",
										Port: &networkingv1alpha3.PortSelector{Number: 9090},
									},
									Weight: 50,
								},
							},
						},
					},
				},
			},
			want: map[string]model.Routes{
				"greeter": []*model.Route{
					{
						ID:             "greeter",
						RequestRetries: 5,
						RequestTimeout: model.DefaultRequestTimeout,
						Source: model.RouteSource{
							Protocol: model.HTTP,
							Port:     8080,
						},
						Targets: []model.RouteTarget{
							{
								Host:    "",
								Port:    9090,
								Weight:  50,
								Version: "v2",
								Type:    model.RouteTargetVersion,
							},
							{
								Host:   "httpbin.myProject.svc.cluster.local",
								Port:   9090,
								Weight: 50,
								Type:   model.RouteTargetExternal,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	i := Istio{kube: kubefake.NewSimpleClientset(), istio: fake.NewSimpleClientset()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := i.istio.NetworkingV1alpha3().VirtualServices(tt.args.projectID).Create(tt.args.ctx, tt.virtualServiceToBeCreated, metav1.CreateOptions{}); err != nil {
				t.Errorf("Cannot generate virutal service required for the test function")
				return
			}
			got, err := i.GetServiceRoutes(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if arr := deep.Equal(got, tt.want); len(arr) > 0 {
				a, _ := json.MarshalIndent(arr, "", " ")
				t.Errorf("GetServiceRoutes() diff = %v", string(a))
			}
			if err := i.istio.NetworkingV1alpha3().VirtualServices(tt.args.projectID).Delete(tt.args.ctx, tt.virtualServiceToBeCreated.ObjectMeta.Name, metav1.DeleteOptions{}); err != nil {
				t.Errorf("Cannot delete virtual service required for the test function")
				return
			}
		})
	}
}
