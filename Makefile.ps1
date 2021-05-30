[string]$OS = "windows"
[string]$ARCH = "amd64"
[string]$SIGNFILES = "publish_files/"

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
    prepare_environment
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
    [array]$OPTIONS = Select-String -Path $CONFIGFILE -Pattern "^(?<key>\w+)=(?<value>.+$)"
    $OPTIONS | ForEach-Object {
        $key, $value = $_.Matches[0].Groups['key','value'].Value
        Set-Variable -Name $key -Value $value
    }
    [string]$BINARY = "$PREFIX-$NAME_v$VERSION"
    [string]$ZIPFILE = "$PREFIX-$NAME_$VERSION_$OS_$ARCH.zip"
    [string]$SHAFILE = "$PREFIX-$NAME_$VERSION_SHA256SUMS"
    Write-Host $PREFIX $NAME $VERSION $BINARY $ZIPFILE $SHAFILE
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
