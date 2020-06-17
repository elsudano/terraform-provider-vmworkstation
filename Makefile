SHELL = /bin/bash

NAME = terraform-provider-vmworkstation
VERSION = v0.1.0
BINARY = $(NAME)_$(VERSION)_x4 
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
	@go build -o $(BINARY)

install: build --copyBIN $(PATHOFPLUGINS)/lock.json ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS)/lock.json
	@sed -i ':N;s/{/{\n  "vmworkstation": "$(shell sha256sum $(BINARY) | awk '{ print $$1 }')",/g' $(PATHOFPLUGINS)/lock.json

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(BINARY)
	@rm -f $(PATHOFPLUGINS)/$(BINARY)
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS)/lock.json


#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--copyBIN:
	@cp $(BINARY) $(PATHOFPLUGINS)/$(BINARY)

.PHONY = $(PHONY)
