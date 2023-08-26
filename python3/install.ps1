param (
  $installationdrive = "C",
  [switch] $createneovimenv = $false
)

if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments
  Break
}

# Python and required pip modules
scoop bucket add versions
scoop install python311
python311 -m pip install --upgrade pip

$currentDIr = Get-Location

if ($createneovimenv -eq $true)
{
  Set-Location "${installationdrive}:\Applications\Scoop\apps\neovim"
  mkdir python3
  Set-Location python3
  mkdir Envs
  Set-Location Envs

  python311 -m venv neovim # Create the virtual env for python3
  .\neovim\Scripts\activate
  py -m pip install pynvim
  py -m pip install neovim
  py -m pip install neovim-remote
  deactivate

  mkdir "~/AppData/Local/python3/Envs"
  New-Item -ItemType SymbolicLink -Path "~/AppData/Local/python3/Envs/neovim" -Target  "${installationdrive}:\Applications\Scoop\apps\neovim\python3\Envs\neovim"
}

Set-Location $currentDir

Read-Host -Prompt "Press Enter to exit"
