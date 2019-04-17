PROJECT_PATH := github.com/spaceuptech/space-cloud
ifeq (,$(wildcard $(HOME)/go/bin))
GOBIN := $(shell go env GOROOT)/bin
else
GOBIN := $(HOME)/go/bin
endif
ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
BINARY_NAME := space-cloud
INSTALL_PATH := /usr/local/bin
SRC_FILES := $(shell find . | grep '.*\.go')

all: $(GOBIN)/linux_amd64/space-cloud $(GOBIN)/darwin_amd64/space-cloud $(GOBIN)/windows_amd64/space-cloud.exe

$(GOBIN)/linux_amd64/space-cloud: *.go $(SRC_FILES)
	GOOS=linux GOARCH=amd64 go install
ifeq (,$(wildcard $(GOBIN)/linux_amd64))
	mkdir $(GOBIN)/linux_amd64
endif
	if [ -f $(GOBIN)/space-cloud ]; then \
		mv $(GOBIN)/space-cloud $(GOBIN)/linux_amd64; \
	fi

$(GOBIN)/darwin_amd64/space-cloud: *.go $(SRC_FILES)
	GOOS=darwin GOARCH=amd64 go install
ifeq (,$(wildcard $(GOBIN)/darwin_amd64))
	mkdir $(GOBIN)/darwin_amd64
endif
	if [ -f $(GOBIN)/space-cloud ]; then \
		mv $(GOBIN)/space-cloud $(GOBIN)/darwin_amd64; \
	fi

$(GOBIN)/windows_amd64/space-cloud.exe: *.go $(SRC_FILES)
	GOOS=windows GOARCH=amd64 go install
ifeq (,$(wildcard $(GOBIN)/windows_amd64))
	mkdir $(GOBIN)/windows_amd64
endif
	if [ -f $(GOBIN)/space-cloud.exe ]; then \
		mv $(GOBIN)/space-cloud.exe $(GOBIN)/windows_amd64; \
	fi

install: $(GOBIN)/$(ARCH)/$(BINARY_NAME)
	install $(GOBIN)/$(ARCH)/$(BINARY_NAME) $(INSTALL_PATH)
	@echo "\x1b[32m$(BINARY_NAME) installed successfully\x1b[0m" >&2

uninstall:
ifeq (,$(wildcard $(INSTALL_PATH)/$(BINARY_NAME)))
	@echo "\x1b[31m$(BINARY_NAME) not installed\x1b[0m" >&2
else
	rm $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "\x1b[32m$(BINARY_NAME) uninstalled successfully\x1b[0m" >&2
endif

zip/: $(GOBIN)/linux_amd64/space-cloud $(GOBIN)/darwin_amd64/space-cloud $(GOBIN)/windows_amd64/space-cloud.exe
ifeq (,$(wildcard zip/))
	mkdir zip
endif
	for file in $$(find $(GOBIN) | grep 'space-cloud'); do \
		outdir=$$(echo $$file | sed -e "s#$(GOBIN)/##" | \
			sed -e 's#/space-cloud.*##') ; \
		outfile=$$(echo $$file | grep -o 'space-cloud.*') ; \
		mkdir zip/$$outdir ; \
		zip zip/$$outdir/$$outfile.zip $$file ; \
	done

zip: zip/

deploy: zip
	curl -H "Authorization: Bearer $(JWT_TOKEN)" \
		-F 'file=@./zip/darwin_amd64/space-cloud.zip' \
		-F 'fileType=file' -F 'makeAll=false' -F 'path=/darwin' \
		https://spaceuptech.com/v1/api/downloads/files
	curl -H "Authorization: Bearer $(JWT_TOKEN)" \
		-F 'file=@./zip/windows_amd64/space-cloud.exe.zip' \
		-F 'fileType=file' -F 'makeAll=false' -F 'path=/windows' \
		https://spaceuptech.com/v1/api/downloads/files
	curl -H "Authorization: Bearer $(JWT_TOKEN)" \
		-F 'file=@./zip/linux_amd64/space-cloud.zip' \
		-F 'fileType=file' -F 'makeAll=false' -F 'path=/linux' \
		https://spaceuptech.com/v1/api/downloads/files

