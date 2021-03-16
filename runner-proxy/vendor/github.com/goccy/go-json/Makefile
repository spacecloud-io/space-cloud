BIN_DIR := $(CURDIR)/bin

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

.PHONY: cover
cover:
	@ go test -coverprofile=cover.tmp.out . ; \
	cat cover.tmp.out | grep -v "encode_optype.go" > cover.out; \
	rm cover.tmp.out

.PHONY: cover-html
cover-html: cover
	go tool cover -html=cover.out

.PHONY: lint
lint: golangci-lint
	golangci-lint run

golangci-lint: | $(BIN_DIR)
	@{ \
		set -e; \
		GOLANGCI_LINT_TMP_DIR=$$(mktemp -d); \
		cd $$GOLANGCI_LINT_TMP_DIR; \
		go mod init tmp; \
		GOBIN=$(BIN_DIR) go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.36.0; \
		rm -rf $$GOLANGCI_LINT_TMP_DIR; \
	}
