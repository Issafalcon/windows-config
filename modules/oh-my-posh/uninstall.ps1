winget uninstall --id JanDeDobbeleer.OhMyPosh -e --source winget --accept-source-agreements

Uninstall-Module -Name Terminal-Icons -Force -ErrorAction SilentlyContinue
# PSReadLine ships with PowerShell — do not uninstall it.

if (Test-Path -Path $PROFILE) {
  Clear-Content $PROFILE -Force
  Write-Host "Cleared $PROFILE"
}
