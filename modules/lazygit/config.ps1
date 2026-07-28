$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$lazygitConfigDir = Join-Path $HOME "AppData\Roaming\lazygit"

if (-not (Test-Path $lazygitConfigDir)) {
  Write-Host "Creating $lazygitConfigDir"
  New-Item -ItemType Directory -Path $lazygitConfigDir | Out-Null
}

$link = Join-Path $lazygitConfigDir "config.yml"
$target = Join-Path $scriptDir "config.yml"
Write-Host "Linking $link -> $target"
if (Test-Path $link) { Remove-Item -Force $link }
try {
  New-Item -ItemType HardLink -Path $link -Target $target | Out-Null
} catch {
  Copy-Item -Force $target $link
  Write-Host "  copied (hardlink unavailable): $_"
}
Write-Host "lazygit config done"
