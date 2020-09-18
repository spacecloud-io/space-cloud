module github.com/spaceuptech/space-cloud/runner

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.4
	github.com/kedacore/keda v1.5.1-0.20200914094616-3a99b77cc330
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/rs/cors v1.7.0
	github.com/segmentio/ksuid v1.0.2
	github.com/spaceuptech/helpers v0.1.2
	github.com/spaceuptech/space-api-go v0.17.3
	github.com/txn2/txeh v1.3.0
	github.com/urfave/cli v1.22.2
	go.etcd.io/bbolt v1.3.3
	gotest.tools v2.2.0+incompatible // indirect; indirects
	istio.io/api v0.0.0-20200518203817-6d29a38039bd
	istio.io/client-go v0.0.0-20200521172153-8555211db875
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.18.8
