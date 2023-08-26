$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)
$omnisharpConfigDir =  "~/.omnisharp"

# Symlink config to correct location
if (Test-Path -Path $omnisharpConfigDir)
{
  "$omnisharpConfigDir exists. Skipping creation"
} else
{
  New-Item -ItemType "directory" -Path $omnisharpConfigDir
}

New-Item -ItemType SymbolicLink -Path "$omnisharpConfigDir/omnisharp.json" -Target "${scriptDir}/omnisharp_formatting.json" -Force
