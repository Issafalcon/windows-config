param ([Parameter(Mandatory)]$installationdrive, $modulename)

# Change the installation location and create if it doesn't exist
[System.Environment]::SetEnvironmentVariable('ChocolateyInstall', "${installationdrive}:\ProgramData\chocolatey", "Machine")
[System.Environment]::SetEnvironmentVariable('ChocolateyToolsLocation', "${installationdrive}:\tools", "Machine")
[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\ProgramData\chocolatey", "Machine")
[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\ProgramData\chocolatey\bin", "Machine")
[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\ProgramData\chocolatey\lib", "Machine")
mkdir "${installationdrive}:\ProgramData\chocolatey"

Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))