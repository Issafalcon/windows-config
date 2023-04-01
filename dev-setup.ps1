###################################################################
# Neovim And Related Dependencies                           
###################################################################



# Node
Write-Host "Installing NodeJs"
choco install nodejs -y

Write-Host "Add nodejs to path"
$env:Path += ";$installationdrive:\Program Files\nodejs" 

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
