.PHONY: pre-hook test-all gen gen-all manifest build run install install-deps

pre-hook:
	go mod tidy
	go mod vendor

gen: pre-hook
	controller-gen object paths="./..."
	./hack/update-codegen.sh

gen-all: manifest gen
	rm -rf vendor

manifest: pre-hook
	controller-gen crd paths="./..." output:crd:artifacts:config=manifests/crds

test-all: build
	go test ./...
	./hack/verify-codegen.sh

build:
	cd cmd/spacectl && go build .

install:
	cd cmd/spacectl && go install

install-deps:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.12.0

run: build
	cd cmd/spacectl && ./spacectl run --kratos-enable