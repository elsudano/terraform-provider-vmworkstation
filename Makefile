SHELL = /bin/bash

CONFIG_FILE = binary_config.ini
ARCH = amd64_v1
SIGNFILES = dist/
PREFIX = $(shell cat $(CONFIG_FILE) | grep "PREFIX=" | tr -d "PREFIX=")
NAME = $(shell cat $(CONFIG_FILE) | grep "NAME=" | tr -d "NAME=")
VERSION = $(shell cat $(CONFIG_FILE) | grep "VERSION=" | tr -d "VERSION=")
# https://semver.org/
BINARY = $(PREFIX)-$(NAME)_v$(VERSION)
SHAFILE = $(PREFIX)-$(NAME)_$(VERSION)_SHA256SUMS

#-------------------------------------------------------#
#    Public Functions                                   #
#-------------------------------------------------------#
PHONY += help
help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

prepare: ## Prepare the environment in order to build the provider
	@go get -u
	@go mod tidy	

build: prepare ## Build the binary of the module
	@git tag v$(VERSION)
	@goreleaser build --clean

publish: ## This option prepare the zip files to publishing in Terraform Registry
	@export GPG_FINGERPRINT=$(shell gpg -k | head -4 | tail -1 | tr -d " ")
	@git add .
	@git commit -m "feat: We have created a new version v$(VERSION)"
	@git tag v$(VERSION)
	@git push
	@gpg --armor --export-secret-keys > private.gpg
	@goreleaser release --clean --skip=publish # --snapshot
	@git push --tags

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@git tag -d v$(VERSION)
	@rm -f $(BINARY) $(BINARY).exe
	@rm -fR $(SIGNFILES)*
	@rm -f private.gpg

.PHONY = $(PHONY)
