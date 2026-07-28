###################################################################
# Oh My Posh + PowerShell profile (CurrentUser — no admin)
###################################################################

winget install JanDeDobbeleer.OhMyPosh -s winget --accept-package-agreements --accept-source-agreements
Install-Module -Name Terminal-Icons -Repository PSGallery -Scope CurrentUser -Force
Install-Module -Name PSReadLine -Scope CurrentUser -Force -SkipPublisherCheck

if (!(Test-Path -Path $PROFILE)) {
  New-Item -ItemType File -Path $PROFILE -Force
}

# MSIX/winget installs no longer set POSH_THEMES_PATH. Pass the theme name
# (no path / .omp.json) so oh-my-posh resolves/downloads it itself.
# Docs: https://ohmyposh.dev/docs/installation/prompt
Clear-Content $PROFILE -Force
@(
  "oh-my-posh init pwsh --config 'easy-term' | Invoke-Expression"
  "Import-Module -Name Terminal-Icons"
  "Set-PSReadLineOption -PredictionSource History"
  "Set-PSReadLineOption -PredictionViewStyle ListView"
  "Set-PSReadLineOption -EditMode Windows"
) | Set-Content -Path $PROFILE -Encoding utf8

Write-Host "Wrote oh-my-posh init to $PROFILE"
