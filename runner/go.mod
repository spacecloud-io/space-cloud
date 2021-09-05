module github.com/spaceuptech/space-cloud/runner

go 1.15

require (
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v8 v8.3.3
	github.com/go-test/deep v1.0.4
	github.com/gogo/protobuf v1.3.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.7.4
	github.com/kedacore/keda v1.5.1-0.20201104115818-50bec808f8b4
	github.com/mitchellh/mapstructure v1.3.3
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/common v0.14.0
	github.com/rs/cors v1.7.0
	github.com/segmentio/ksuid v1.0.2
	github.com/spaceuptech/helpers v0.2.1
	github.com/spaceuptech/space-api-go v0.17.3
	github.com/urfave/cli v1.22.2
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
	istio.io/api v0.0.0-20200518203817-6d29a38039bd
	istio.io/client-go v0.0.0-20200521172153-8555211db875
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.18.8
