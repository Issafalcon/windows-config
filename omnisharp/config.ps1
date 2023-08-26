$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)
$omnisharpConfigDir =  "~/.omnisharp"

if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments
  Break
}

# Symlink config to correct location
if (Test-Path -Path $omnisharpConfigDir)
{
  "$omnisharpConfigDir exists. Skipping creation"
} else
{
  New-Item -ItemType "directory" -Path $omnisharpConfigDir
}

New-Item -ItemType SymbolicLink -Path "$omnisharpConfigDir/omnisharp.json" -Target "${scriptDir}/omnisharp_formatting.json" -Force

Read-Host -Prompt "Press Enter to exit"
