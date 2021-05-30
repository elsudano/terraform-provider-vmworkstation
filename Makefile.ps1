[string]$CONFIGFILE = ".\binary_config.ini"
[array]$CONFIGURATIONS = Get-Content $CONFIGFILE

foreach ($CONFI in $CONFIGURATIONS) {
    [string]$name, [string]$value = $CONFI -split "=", 2
    Write-Host "Valor de la variable $name : $value"
}
# $PREFIX = terraform-provider
# $NAME = vmworkstation
# $VERSION = 0.2.2