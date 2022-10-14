.PHONY: pre-hook generate manifest build

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

build:
	cd cmd/spacectl && go build .

run: build
	cd cmd/spacectl && ./spacectl run