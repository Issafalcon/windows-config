param (
  $installationdrive = "C"
)

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)
New-Item -ItemType SymbolicLink -Path "~/.config/wezterm" -Target $scriptDir/.config/