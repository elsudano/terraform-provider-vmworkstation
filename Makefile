SHELL = /bin/bash

PREFIX = terraform-provider
NAME = vmworkstation
VERSION = 0.1.0
OS = linux
ARCH = amd64
DIRELEASES = 
BINARY = $(PREFIX)-$(NAME)_v$(VERSION)
ZIPFILE = $(PREFIX)-$(NAME)_$(VERSION)_$(OS)_$(ARCH).zip
SHAFILE = $(PREFIX)-$(NAME)_$(VERSION)_SHA256SUMS
PATHOFTERRAFORM = ~/GitHub/Automated-Deploy-Server/.terraform
PATHOFPLUGINS_12 = $(PATHOFTERRAFORM)/plugins/linux_amd64
PATHOFPLUGINS_13 = $(PATHOFTERRAFORM)/plugins/registry.terraform.io/elsudano/$(NAME)/$(VERSION)/$(OS)_$(ARCH)

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

publish: build --compress ## This option prepare the zip files to publishing in Terraform Registry
	@sha256sum $(DIRELEASES)*.zip > $(DIRELEASES)$(SHAFILE)
	@gpg -q --detach-sign $(DIRELEASES)$(SHAFILE)

install: build --moveBIN $(PATHOFPLUGINS_12)/lock.json ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS_12)/lock.json
	@sed -i ':N;s/{/{\n  "vmworkstation": "$(shell sha256sum $(PATHOFPLUGINS_12)/$(BINARY)_x4 | awk '{ print $$1 }')",/g' $(PATHOFPLUGINS_12)/lock.json
	@ls -lah $(PATHOFPLUGINS_12) $(PATHOFPLUGINS_13)

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(DIRELEASES)$(BINARY) $(DIRELEASES)$(ZIPFILE) $(DIRELEASES)$(SHAFILE) $(DIRELEASES)$(SHAFILE).sig
	@rm -fR $(PATHOFTERRAFORM)

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--moveBIN:
	@mkdir --parents $(PATHOFPLUGINS_12)
	@mkdir --parents $(PATHOFPLUGINS_13)
	@cp $(DIRELEASES)$(BINARY) $(PATHOFPLUGINS_12)/$(BINARY)_x4
	@mv $(DIRELEASES)$(BINARY) $(PATHOFPLUGINS_13)/$(BINARY)

--compress:
	@zip -q $(DIRELEASES)$(ZIPFILE) $(DIRELEASES)$(BINARY)

.PHONY = $(PHONY)
