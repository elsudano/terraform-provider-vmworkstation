Function help {
    [array]$OPTIONS = Get-Content ".\Makefile.ps1"
    foreach ($CONFI in $CONFIGURATIONS) {
        Select-String -Pattern "^Function\w?<fun>?{"
        [string]$name, [string]$value = $CONFI -split "##", 2
        Write-Host "Valor de la variable $name : $value"
    }
}

Function clean { ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
	# @rm -f $(BINARY)
	# @rm -f $(SIGNFILES)/*
}

Function compress { ## With this function we comppress the files in one, and the we calculate the sha256
	# @zip -q $(SIGNFILES)$(ZIPFILE) $(BINARY)
}

Function build { ## Build the binary of the module

}

Function install { ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
    build
}

Function publish { ## This option prepare the zip files to publishing in Terraform Registry
    clean
    install 
    compress
}

Function prepare_environment {
    [string]$CONFIGFILE = ".\binary_config.ini"
    [array]$CONFIGURATIONS = Get-Content $CONFIGFILE

    foreach ($CONFI in $CONFIGURATIONS) {
        [string]$name, [string]$value = $CONFI -split "=", 2
        Write-Host "Valor de la variable $name : $value"
    }
    # $PREFIX = terraform-provider
    # $NAME = vmworkstation
    # $VERSION = 0.2.2
}
help