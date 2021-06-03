[string]${OS} = "windows"
[string]${ARCH} = "amd64"
[string]$SIGNFILES = "publish_files/"

[string]${global:VERSION}
[string]${global:BINARY}
[string]${global:ZIPFILE}
[string]${global:SHAFILE}

Function help {
    [array]${HELPS} = Select-String -Path .\Makefile.ps1 -Pattern "^(Function?) (?<fun>\w+) (.+##) (?<help>.+)"
    ${HELPS} | ForEach-Object {
        ${var}, ${help} = $_.Matches[0].Groups['fun','help'].Value
        Write-Host "${var} : `t${help}"
    }
}

Function clean { ## Clean the project, this only remove default config of API REST VmWare Workstation Pro, the cert, private key and binary
    if ( Test-Path -Path ${global:BINARY} ) {
        Remove-Item ${global:BINARY}
    }
    if ( Test-Path -Path ${global:SHAFILE} ) {
        Remove-Item ${global:SHAFILE}
    }
    if ( Test-Path -Path ${SIGNFILES} ) {
        Remove-Item ${SIGNFILES} -Recurse
    }
    Write-Host "we have did Cleaning"
}
 
Function prepare_environment {
    [string]${CONFIGFILE} = ".\binary_config.ini"
    [array]${OPTIONS} = Select-String -Path $CONFIGFILE -Pattern "^(?<key>\w+)=(?<value>.+$)"
    ${OPTIONS} | ForEach-Object {
        ${key}, ${value} = $_.Matches[0].Groups['key','value'].Value
        Set-Variable -Name ${key} -Value ${value}
    }
    ${global:VERSION} = ${VERSION}
    ${global:BINARY} = "${PREFIX}-${NAME}_v${VERSION}.exe"
    ${global:ZIPFILE} = "${PREFIX}-${NAME}_${VERSION}_${OS}_${ARCH}.zip"
    ${global:SHAFILE} = "${PREFIX}-${NAME}_${VERSION}_SHA256SUMS"
    #Write-Host ${PREFIX} ${NAME} ${VERSION} ${global:BINARY} ${global:ZIPFILE} ${global:SHAFILE}
}

Function build { ## Build the binary of the module
    prepare_environment
    & go build -o ${global:BINARY}
    Write-Host "we made the binary"
}

Function compress { ## With this function we comppress the files in one, and the we calculate the sha256
    build
    if ( -Not (Test-Path -Path ${SIGNFILES}) ) {
        New-Item ${SIGNFILES} -itemtype directory
    }
    Compress-Archive -path ${global:BINARY} -destinationpath ${SIGNFILES}${global:ZIPFILE}
    Write-Host "we have did Compressing"
}

Function install { ## Copy binary to the project and det SHA256SUM in the config of project, NOTE: Just for Dev. environment for both Terraform 0.12 and 0.13_beta2
    build
    [string]${PLUGIN_PATH} = "$env:APPDATA\terraform.d\plugins\registry.terraform.io\elsudano\vmworkstation\${global:VERSION}\${OS}_${ARCH}\"
    if ( -Not (Test-Path -Path ${PLUGIN_PATH}) ) {
        New-Item ${PLUGIN_PATH} -itemtype directory
    }
    Copy-Item -Path ${global:BINARY} -Destination ${PLUGIN_PATH}
    Write-Host "When you to be develop a provider, is better use the ~/.terraformrc file"
}

Function publish { ## This option prepare the zip files to publishing in Terraform Registry
    clean 
    compress
    [string]${HASH} = Get-FileHash -Path ${SIGNFILES}${global:ZIPFILE} | Select-Object Hash
    ${HASH} = ${HASH} -replace '@{Hash=','' 
    ${HASH} = ${HASH} -replace '}','' 
    ${HASH} | Out-File -FilePath ${global:SHAFILE}
    Write-Host "we have did Publish and the Hash is: ${HASH}"
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
