param (
  [Parameter(Mandatory)]$installationdrive,
  [switch] $createneovimenv = $true
)

# Python and required pip modules
scoop bucket add versions
scoop install python311
python3 -m pip install --upgrade pip

$currentDIr = Get-Location
Set-Location "${installationdrive}:\Applications\Scoop\apps\neovim"
mkdir python3
Set-Location python3
mkdir Envs
Set-Location Envs

if ($createneovimenv -eq $true) {
  python3 -m venv neovim # Create the virtual env for python3
  .\neovim\Scripts\activate
  python3 -m pip install pynvim
  python3 -m pip install neovim
  python3 -m pip install neovim-remote
  deactivate

  mkdir "~/AppData/Local/python3/Envs"
  New-Item -ItemType SymbolicLink -Path "~/AppData/Local/python3/Envs/neovim" -Target  "${installationdrive}:\Applications\Scoop\apps\neovim\python3\Envs\neovim"
}

Set-Location $currentDir
