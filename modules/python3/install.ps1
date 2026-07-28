param (
  $installationdrive = "C",
  [switch] $createneovimenv = $false
)

# Scoop must run unelevated (admin sessions abort / break user-scoped installs).
scoop bucket add versions
scoop install python311
python311 -m pip install --upgrade pip

$currentDir = Get-Location

if ($createneovimenv) {
  $venvRoot = "${installationdrive}:\Applications\Scoop\apps\neovim\python3"
  $venvPath = Join-Path $venvRoot "Envs\neovim"
  New-Item -ItemType Directory -Force -Path (Join-Path $venvRoot "Envs") | Out-Null

  python311 -m venv $venvPath
  & (Join-Path $venvPath "Scripts\Activate.ps1")
  python -m pip install pynvim neovim neovim-remote
  deactivate

  # Directory junction: no admin / Developer Mode required (unlike SymbolicLink).
  $linkParent = Join-Path $HOME "AppData\Local\python3\Envs"
  $linkPath = Join-Path $linkParent "neovim"
  New-Item -ItemType Directory -Force -Path $linkParent | Out-Null
  if (Test-Path $linkPath) {
    Remove-Item -Force -Recurse $linkPath
  }
  New-Item -ItemType Junction -Path $linkPath -Target $venvPath | Out-Null
}

Set-Location $currentDir
