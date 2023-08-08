param (
  $installationdrive = "C"
)

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

# Symlink config to correct location
New-Item -ItemType SymbolicLink -Path "~/.ideavimrc" -Target "${scriptDir}/.ideavimrc" -Force
New-Item -ItemType SymbolicLink -Path "~/.vimrc" -Target "${scriptDir}/.vimrc" -Force

if ($installationdrive -ne "C")
{
  New-Item -ItemType SymbolicLink -Path "C:/Program Files/Git/" -Target "${installationdrive}:/Program Files/Git/"
}
