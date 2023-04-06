$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

# Symlink config to correct location
New-Item -ItemType "directory" -Path "~/AppData/Roaming/lazygit"
New-Item -ItemType SymbolicLink -Path "~/AppData/Roaming/lazygit/config.yml" -Target "${scriptDir}/config.yml"
