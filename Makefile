SHELL = /bin/bash

PREFIX = terraform-provider
NAME = vmworkstation
VERSION = 0.2.2
# https://semver.org/
OS = linux
ARCH = amd64
SIGNFILES = publish_files/
BINARY = $(PREFIX)-$(NAME)_v$(VERSION)
ZIPFILE = $(PREFIX)-$(NAME)_$(VERSION)_$(OS)_$(ARCH).zip
SHAFILE = $(PREFIX)-$(NAME)_$(VERSION)_SHA256SUMS

#-------------------------------------------------------#
#    Public Functions                                   #
#-------------------------------------------------------#
PHONY += help
help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary of the module
	@go build -o $(BINARY)

install: build ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
	@echo When you to be develop a provider, is better use the ~/.terraformrc file
	@cat ~/.terraformrc | grep -B 5 -A 2 $(NAME)
	@ls -lah $(BINARY)

publish: clean install --compress ## This option prepare the zip files to publishing in Terraform Registry
	@cd $(SIGNFILES); sha256sum *.zip > $(SHAFILE)
	@cd $(SIGNFILES); gpg -q --detach-sign $(SHAFILE)

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(BINARY)
	@rm -f $(SIGNFILES)/*

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--compress:
	@zip -q $(SIGNFILES)$(ZIPFILE) $(BINARY)

.PHONY = $(PHONY)
