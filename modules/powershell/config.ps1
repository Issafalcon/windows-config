###################################################################
# Windows Powershell (CurrentUser — no admin)
###################################################################

winget install JanDeDobbeleer.OhMyPosh -s winget
Install-Module -Name Terminal-Icons -Repository PSGallery -Scope CurrentUser -Force
Install-Module -Name PSReadLine -Scope CurrentUser -Force -SkipPublisherCheck

if (!(Test-Path -Path $PROFILE)) {
  New-Item -ItemType File -Path $PROFILE -Force
}

Clear-Content $PROFILE -Force
Add-Content -Value "oh-my-posh --init --shell pwsh --config `$env:POSH_THEMES_PATH`\easy-term.omp.json | Invoke-Expression" -Path $PROFILE
Add-Content -Value "Import-Module -Name Terminal-Icons" -Path $PROFILE
Add-Content -Value "Set-PSReadLineOption -PredictionSource History" -Path $PROFILE
Add-Content -Value "Set-PSReadLineOption -PredictionViewStyle ListView" -Path $PROFILE
Add-Content -Value "Set-PSReadLineOption -EditMode Windows" -Path $PROFILE
