Write-Host "Installing Choco"

Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

refreshenv

###################################################################
# Visual Studio
###################################################################
choco install visualstudio2022community # Upgrade to professional or enterprise as needed
choco install visualstudio2022-workload-azure
choco install visualstudio2022-workload-netweb
choco install visualstudio2022-workload-netcoretools
choco install visualstudio2022-workload-nativedesktop # Needed even if not usign C++ dev as the netcoredbg install relies on it

###################################################################
# Neovim And Related Dependencies                           
###################################################################

# Neovim
choco install neovim -y

# Git
winget install --id Git.Git -e --source winget
choco install delta -y

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
