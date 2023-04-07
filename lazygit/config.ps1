$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)
$lazygitConfigDir =  "~/AppData/Roaming/lazygit"

# Symlink config to correct location
if (Test-Path -Path $lazygitConfigDir) {
    "$lazygitConfigDir exists. Skipping creation"
} else {
  New-Item -ItemType "directory" -Path "~/AppData/Roaming/lazygit"
}

New-Item -ItemType SymbolicLink -Path "~/AppData/Roaming/lazygit/config.yml" -Target "${scriptDir}/config.yml" -Force
