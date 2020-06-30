SHELL = /bin/bash

PREFIX = terraform-provider
NAME = vmware-workstation
VERSION = 0.1.0
OS = linux
ARCH = x86_64
DIRELEASES = releases
BINARY = $(PREFIX)-$(NAME)_v$(VERSION)
ZIPFILE = $(PREFIX)-$(NAME)_$(VERSION)_$(OS)_$(ARCH).zip
SHAFILE = $(PREFIX)-$(NAME)_$(VERSION)_SHA256SUMS
PATHOFPLUGINS = ~/GitHub/Automated-Deploy-Server/.terraform/plugins/linux_amd64

#-------------------------------------------------------#
#    Public Functions                                   #
#-------------------------------------------------------#
PHONY += help
help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary of the module
	@go build -o $(DIRELEASES)/$(BINARY)

publish: build --compress ## This option prepare the zip files to publishing in Terraform Registry
	@sha256sum $(DIRELEASES)/*.zip > $(DIRELEASES)/$(SHAFILE)
	@gpg --detach-sign $(DIRELEASES)/$(SHAFILE)

install: build --copyBIN $(PATHOFPLUGINS)/lock.json ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS)/lock.json
	@sed -i ':N;s/{/{\n  "vmworkstation": "$(shell sha256sum $(BINARY) | awk '{ print $$1 }')",/g' $(PATHOFPLUGINS)/lock.json

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(DIRELEASES)/$(BINARY) $(DIRELEASES)/$(ZIPFILE) $(DIRELEASES)/$(SHAFILE) $(DIRELEASES)/$(SHAFILE).sig
	@rm -f $(PATHOFPLUGINS)/$(BINARY)_x4
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS)/lock.json

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--copyBIN:
	@cp $(DIRELEASES)/$(BINARY) $(PATHOFPLUGINS)/$(BINARY)_x4

--compress:
	@zip -q $(DIRELEASES)/$(ZIPFILE) $(DIRELEASES)/$(BINARY)

.PHONY = $(PHONY)
