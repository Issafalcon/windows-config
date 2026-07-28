$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$lazygitConfigDir = Join-Path $HOME "AppData\Roaming\lazygit"

if (-not (Test-Path $lazygitConfigDir)) {
  New-Item -ItemType Directory -Path $lazygitConfigDir | Out-Null
}

$link = Join-Path $lazygitConfigDir "config.yml"
$target = Join-Path $scriptDir "config.yml"
if (Test-Path $link) { Remove-Item -Force $link }
try {
  New-Item -ItemType HardLink -Path $link -Target $target | Out-Null
} catch {
  Copy-Item -Force $target $link
}
