SHELL = /bin/bash

NAME = terraform-provider-vmworkstation
VERSION = v0.0.1
BINARY = $(NAME)_$(VERSION)_x4 
PRIVATEKEYFILE = workstationapi-key.pem
CERTFILE = workstationapi-cert.pem
IPADDRESS = 127.0.0.1
PORT = 5555
PATHOFPLUGINS = ~/GitHub/Automated-Deploy-Server/.terraform/plugins/linux_amd64

#-------------------------------------------------------#
#    Public Functions                                   #
#-------------------------------------------------------#
PHONY += help
help:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

PHONY += bootstrap 
bootstrap: --generateSSL --vmrest ## Prepare environment for you can use a API REST of VmWare Workstation Pro and generate files for SSL

build: ## Build the binary of the module
	@go build -o $(BINARY)

install: build --copyBIN $(PATHOFPLUGINS)/lock.json ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS)/lock.json
	@sed -i ':N;s/{/{\n  "vmworkstation": "$(shell sha256sum $(BINARY) | awk '{ print $$1 }')",/g' $(PATHOFPLUGINS)/lock.json

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@rm -f $(PRIVATEKEYFILE)
	@rm -f $(CERTFILE)
	@rm -f $(BINARY)
	@rm -f $(PATHOFPLUGINS)/$(BINARY)
	@sed -i '/"vmworkstation":/'d $(PATHOFPLUGINS)/lock.json
	@rm -f ~/.vmrestCfg

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--generateSSL:
	@openssl req -x509 -newkey rsa:4096 -keyout $(PRIVATEKEYFILE) -out $(CERTFILE) -days 365 -nodes -subj "/C=ES/ST=Granada/L=Granada/O=Internet SL/OU=IT/CN=localhost"
	@vmrest -C

--copyBIN:
	@cp $(BINARY) $(PATHOFPLUGINS)/$(BINARY)

--vmrest: $(PRIVATEKEYFILE) $(CERTFILE)
	@vmrest -k $(PRIVATEKEYFILE) -c $(CERTFILE) -i $(IPADDRESS) -p $(PORT)

.PHONY = $(PHONY)
