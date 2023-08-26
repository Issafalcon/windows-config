param (
  $installationdrive = "C"
)

if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments
  Break
}

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

# Symlink config to correct location
New-Item -ItemType SymbolicLink -Path "~/.ideavimrc" -Target "${scriptDir}/.ideavimrc" -Force
New-Item -ItemType SymbolicLink -Path "~/.vimrc" -Target "${scriptDir}/.vimrc" -Force

Read-Host -Prompt "Press Enter to exit"
