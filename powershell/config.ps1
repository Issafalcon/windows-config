###################################################################
# Windows Powershell
###################################################################

if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process pwsh -Verb runAs -ArgumentList $arguments -Wait
  Break
}

winget install JanDeDobbeleer.OhMyPosh -s winget
Install-Module -Name Terminal-Icons -Repository PSGallery
Install-Module -Name PSReadLine -Scope AllUsers -Force -SkipPublisherCheck
if (!(Test-Path -Path $PROFILE.AllUsersAllHosts)) {
  New-Item -ItemType File -Path $PROFILE.AllUsersAllHosts -Force
}

Clear-Content $PROFILE.AllUsersAllHosts -Force
Add-Content -Value "oh-my-posh --init --shell pwsh --config `$env:POSH_THEMES_PATH`\easy-term.omp.json | Invoke-Expression" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Import-Module -Name Terminal-Icons" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Set-PSReadLineOption -PredictionSource History" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Set-PSReadLineOption -PredictionViewStyle ListView" -Path $PROFILE.AllUsersAllHosts
Add-Content -Value "Set-PSReadLineOption -EditMode Windows" -Path $PROFILE.AllUsersAllHosts

Read-Host -Prompt "Press Enter to exit"
