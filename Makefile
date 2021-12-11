SHELL = /bin/bash

CONFIG_FILE = binary_config.ini
ARCH = amd64
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

build: ## Build the binary of the module
	@git tag v$(VERSION)
	@goreleaser build --rm-dist

install: build ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
	@echo When you to be develop a provider, is better use the ~/.terraformrc file
	@cat ~/.terraformrc | grep -B 5 -A 2 $(NAME)
	@ls -lahr $(SIGNFILES)

publish: install --compress ## This option prepare the zip files to publishing in Terraform Registry
#@export GPG_FINGERPRINT=$(shell gpg --armor --export-secret-keys | tr -d \\n);
	@goreleaser release --rm-dist # --snapshot
#@cd $(SIGNFILES); sha256sum *.zip > $(SHAFILE)
#@cd $(SIGNFILES); gpg -q --detach-sign $(SHAFILE)

clean: ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	@git tag -d v$(VERSION)
	@rm -f $(BINARY) $(BINARY).exe
	@rm -fR $(SIGNFILES)*

#-------------------------------------------------------#
#    Private Functions                                  #
#-------------------------------------------------------#
--compress:
	@zip -q $(SIGNFILES)$(PREFIX)-$(NAME)_$(VERSION)_linux_$(ARCH).zip $(SIGNFILES)$(PREFIX)-$(NAME)_linux_$(ARCH)/$(BINARY)
	@zip -q $(SIGNFILES)$(PREFIX)-$(NAME)_$(VERSION)_windows_$(ARCH).zip $(SIGNFILES)$(PREFIX)-$(NAME)_windows_$(ARCH)/$(BINARY).exe

.PHONY = $(PHONY)
