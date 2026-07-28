param (
  $installationdrive = "C"
)

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

function Set-FileLink([string]$Link, [string]$Target) {
  if (Test-Path $Link) { Remove-Item -Force $Link }
  try {
    New-Item -ItemType HardLink -Path $Link -Target $Target | Out-Null
  } catch {
    Copy-Item -Force $Target $Link
  }
}

Set-FileLink (Join-Path $HOME ".ideavimrc") (Join-Path $scriptDir ".ideavimrc")
Set-FileLink (Join-Path $HOME ".vimrc") (Join-Path $scriptDir ".vimrc")
