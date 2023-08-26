if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments
  Break
}

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)
$lazygitConfigDir =  "~/AppData/Roaming/lazygit"

# Symlink config to correct location
if (Test-Path -Path $lazygitConfigDir) {
    "$lazygitConfigDir exists. Skipping creation"
} else {
  New-Item -ItemType "directory" -Path "~/AppData/Roaming/lazygit"
}

New-Item -ItemType SymbolicLink -Path "~/AppData/Roaming/lazygit/config.yml" -Target "${scriptDir}/config.yml" -Force
