SHELL = /bin/bash

PREFIX = terraform-provider
NAME = vmworkstation
VERSION = 0.2.0
# https://semver.org/
OS = linux
ARCH = amd64
DIRELEASES = releases/
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
	@go build -o $(DIRELEASES)$(BINARY)

install: build ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
	@echo When you to be develop a provider, is better use the ~/.terraformrc file
	@cat ~/.terraformrc | grep -C 3 $(NAME)
	@ls -lah $(DIRELEASES)

publish: clean install --compress ## This option prepare the zip files to publishing in Terraform Registry
	@sha256sum $(DIRELEASES)*.zip > $(DIRELEASES)$(SHAFILE)
	@gpg -q --detach-sign $(DIRELEASES)$(SHAFILE)
	@mv $(DIRELEASES)$(ZIPFILE) $(SIGNFILES)
	@mv $(DIRELEASES)$(SHAFILE) $(SIGNFILES)
	@mv $(DIRELEASES)$(SHAFILE).sig $(SIGNFILES)

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(DIRELEASES)$(BINARY)
	@rm -f $(SIGNFILES)/*

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--compress:
	@zip -q $(DIRELEASES)$(ZIPFILE) $(DIRELEASES)$(BINARY)

.PHONY = $(PHONY)
