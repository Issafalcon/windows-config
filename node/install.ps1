param (
  [string]$installationdrive = "C"
)

$url = "https://github.com/coreybutler/nvm-windows/releases/download/1.1.12/nvm-setup.exe"
Invoke-WebRequest -Uri $url -OutFile $PSScriptRoot\nvm-setup.exe

$args = @("Comma", "Separated", "Arguments")
Start-Process -Filepath "$PSScriptRoot/nvm-setup.exe" -ArgumentList $args -Wait

Remove-Item $PSScriptRoot/nvm-setup.exe -Force

[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}\Program Files\nodejs", "Machine")
nvm install lts
nvm use lts
