$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$omnisharpConfigDir = Join-Path $HOME ".omnisharp"

if (-not (Test-Path $omnisharpConfigDir)) {
  New-Item -ItemType Directory -Path $omnisharpConfigDir | Out-Null
}

$link = Join-Path $omnisharpConfigDir "omnisharp.json"
$target = Join-Path $scriptDir "omnisharp_formatting.json"
if (Test-Path $link) { Remove-Item -Force $link }
try {
  New-Item -ItemType HardLink -Path $link -Target $target | Out-Null
} catch {
  Copy-Item -Force $target $link
}
