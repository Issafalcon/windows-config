param (
  [Parameter(Mandatory)] $installationdrive = "default"
)

Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod get.scoop.sh -outfile 'install.ps1'

if ($installationdrive -eq "default") {
  $installationdrive = "C"
}

& .\install.ps1 -ScoopDir "${installationdrive}:\Applications\Scoop" -ScoopGlobalDir "${installationdrive}:\GlobalScoopApps" -NoProxy

Remove-Item .\install.ps1
