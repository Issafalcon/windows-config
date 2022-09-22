###################################################################
# Windows Powershell
###################################################################

winget install JanDeDobbeleer.OhMyPosh -s winget
Install-Module -Name Terminal-Icons -Repository PSGallery
if (!(Test-Path -Path $PROFILE.AllUsersAllHosts)) {
  New-Item -ItemType File -Path $PROFILE.AllUsersAllHosts -Force
}
Add-Content -Value "oh-my-posh --init --shell pwsh --config ~/jandedobbeleer.omp.json | Invoke-Expression" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Import-Module -Name Terminal-Icons" -Path $PROFILE.AllUsersAllHosts
refreshenv
