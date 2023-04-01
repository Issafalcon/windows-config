param ([Parameter(Mandatory)]$installationdrive, $installvisualstudio)
Write-Host "Installing Choco"

$env:ChocolateyInstall = '$installationdrive:\ProgramData\chocolatey\'
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

refreshenv

###################################################################
# Neovim And Related Dependencies                           
###################################################################

# Neovim
choco install neovim -y

# Install gcc (also adds tar)
choco install mingw
$env:Path += ';C:\ProgramData\chocolatey\lib\mingw\tools\install\mingw64\bin'

# 7zip
choco install 7zip -y

# Python and required pip modules
choco install python3 -y
cd ~
mkdir Envs
cd Envs
py -m venv neovim3 # Create the virtual env for python3
.\neovim3\Scripts\activate
py -m pip install pynvim
py -m pip install neovim
py -m pip install neovim-remote
deactivate

# Go
choco install go -y

# RipGrep, fd, fzf and Silver Searcher for finding and sorting
choco install ripgrep -y
choco install ag -y
choco install fzf -y
choco install fd -y

# Node
Write-Host "Installing NodeJs"
choco install nodejs -y

Write-Host "Add nodejs to path"
$env:Path += ";C:\Program Files\nodejs" 

# Install neovim
npm install -g neovim
#
## Cmake
#New-Item -ItemType Directory -Path "C:\temp"
#Invoke-WebRequest -Uri "https://github.com/Kitware/CMake/releases/download/v3.24.0-rc3/cmake-3.24.0-rc3-windows-x86_64.msi" -OutFile "C:\temp\cmake.msi"
#Start-Process msiexec -Wait -ArgumentList "/I C:\temp\cmake.msi"
#$env:Path += ';C:\Program Files\CMake\bin'
#
#
## wget
#choco install wget -y
#
## gzip
#choco install gzip -y
#py -m pip install --upgrade pip
#
## Tree Sitter
#choco install tree-sitter -y
#
## If issue https://github.com/Samsung/netcoredbg/issues/94 still not fixed, may have to manually modify the source code prior to this command
#cmake .. -G "Visual Studio 17 2022" # Assumes we have 2022 VS installed
refreshenv
