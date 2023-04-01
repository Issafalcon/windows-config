param ([Parameter(Mandatory)]$installationdrive, $installvisualstudio)

# Git
winget install --id Git.Git -e --source winget --location "${installationdrive}:\Program Files\Git"
[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\Program Files\Git\cmd", "Machine")
[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\ProgramData\chocolatey\lib\delta\tools", "Machine")
choco install delta -y
