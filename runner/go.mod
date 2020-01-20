module github.com/spaceuptech/space-cloud/runner

go 1.13

require (
	github.com/containerd/containerd v1.3.2 // indirect
	github.com/coreos/bbolt v1.3.3
	github.com/dgraph-io/badger v1.6.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190927142053-ada3c14355ce
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/rs/cors v1.7.0
	github.com/segmentio/ksuid v1.0.2
	github.com/sirupsen/logrus v1.4.2
	github.com/txn2/txeh v1.3.0
	github.com/urfave/cli v1.22.2
	go.etcd.io/bbolt v1.3.3
	istio.io/api v0.0.0-20191109011911-e51134872853
	istio.io/client-go v0.0.0-20191206191348-5c576a7ecef0
	k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
)
