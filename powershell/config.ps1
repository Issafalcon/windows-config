###################################################################
# Windows Powershell
###################################################################

winget install JanDeDobbeleer.OhMyPosh -s winget
Install-Module -Name Terminal-Icons -Repository PSGallery
Install-Module -Name PSReadLine -AllowPrerelease -Scope AllUsers -Force -SkipPublisherCheck
if (!(Test-Path -Path $PROFILE.AllUsersAllHosts)) {
  New-Item -ItemType File -Path $PROFILE.AllUsersAllHosts -Force
}
Add-Content -Value "oh-my-posh --init --shell pwsh --config $env:POSH_THEMES_PATH\easy-term.omp.json | Invoke-Expression" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Import-Module -Name Terminal-Icons" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Set-PSReadLineOption -PredictionSource History" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Set-PSReadLineOption -PredictionViewStyle ListView" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Set-PSReadLineOption -EditMode Windows" -Path $PROFILE.AllUsersAllHosts
refreshenv
