Function help {
    [array]$HELPS = Select-String -Path .\Makefile.ps1 -Pattern "^(Function?) (?<fun>\w+) (.+##) (?<help>.+)"
    $HELPS | ForEach-Object {
        $var, $help = $_.Matches[0].Groups['fun','help'].Value
        Write-Host "$var : `t$help"
    }
}

Function clean { ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
    Write-Host "we did Clean"
}

Function compress { ## With this function we comppress the files in one, and the we calculate the sha256
    Write-Host "we did Compress"
}

Function build { ## Build the binary of the module
    Write-Host "we did Build"
}

Function install { ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
    build
    Write-Host "we did Install"
}

Function publish { ## This option prepare the zip files to publishing in Terraform Registry
    clean
    install 
    compress
    Write-Host "we did Publish"
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
# --------------------------
#       Menu
# --------------------------
if ($args.Count -eq 1 ) {
    switch ( $args[0] ) {
        clean { clean }
        compress { compress }
        build { build }
        install { install }
        publish { publish }
    }
} else {
    help
}
