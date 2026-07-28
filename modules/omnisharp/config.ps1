$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$omnisharpConfigDir = Join-Path $HOME ".omnisharp"

if (-not (Test-Path $omnisharpConfigDir)) {
  Write-Host "Creating $omnisharpConfigDir"
  New-Item -ItemType Directory -Path $omnisharpConfigDir | Out-Null
}

$link = Join-Path $omnisharpConfigDir "omnisharp.json"
$target = Join-Path $scriptDir "omnisharp_formatting.json"
Write-Host "Linking $link -> $target"
if (Test-Path $link) { Remove-Item -Force $link }
try {
  New-Item -ItemType HardLink -Path $link -Target $target | Out-Null
} catch {
  Copy-Item -Force $target $link
  Write-Host "  copied (hardlink unavailable): $_"
}
Write-Host "omnisharp config done"
