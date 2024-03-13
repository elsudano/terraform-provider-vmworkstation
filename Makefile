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
	@export GPG_FINGERPRINT=$(shell gpg -k | head -4 | tail -1 | tr -d " ")
	@export GOPRIVATE=github.com/elsudano/vmware-workstation-api-client; go get github.com/elsudano/vmware-workstation-api-client@$(shell git -C ../vmware-workstation-api-client/ tag --sort=committerdate | tail -1)
	@go get -u
	@go mod tidy	

build: prepare ## Build the binary of the module
	@git add .
	@git commit -m "update: We have updated dependencies before to build"
	@git tag v$(VERSION)
	@goreleaser build --clean

install: build ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
	@echo When you are developing a provider, is better use the ~/.terraformrc file
	@echo "Before to publish you need run these commands:"
	@#cat ~/.terraformrc | grep -B 2 -A 2 $(NAME)
	@#ls -lahr $(SIGNFILES)


publish: install ## This option prepare the zip files to publishing in Terraform Registry
	@go get -u
	@go mod tidy
	@git add .
	@git commit -m "chore: We have updated the dependencies"
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
