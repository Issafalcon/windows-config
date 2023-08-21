param (
  [string]$installationdrive = "C"
)

$url = "https://github.com/coreybutler/nvm-windows/releases/download/1.1.11/nvm-setup.exe"
Invoke-WebRequest -Uri $url -OutFile $PSScriptRoot\nvm-setup.exe

$args = @("Comma", "Separated", "Arguments")
Start-Process -Filepath "$PSScriptRoot/nvm-setup.exe" -ArgumentList $args

Remove-Item $PSScriptRoot/nvm-setup.exe