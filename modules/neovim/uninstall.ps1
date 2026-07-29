# Reverse packages / links this module installs. Shared deps (gh, python3, node)
# stay — uninstall those modules separately.
npm uninstall -g tree-sitter-cli 2>$null

foreach ($pkg in @(
  'win32yank', 'gzip', 'unzip', 'wget',
  'fd', 'fzf', 'ag', 'ripgrep',
  'neovim', 'vcredist2022'
)) {
  scoop uninstall $pkg 2>$null
}

# Leave mingw/go — often shared; clear treesitter compiler hints only.
[System.Environment]::SetEnvironmentVariable('CC', $null, 'User')
[System.Environment]::SetEnvironmentVariable('CXX', $null, 'User')
[System.Environment]::SetEnvironmentVariable('EDITOR', $null, 'User')
[System.Environment]::SetEnvironmentVariable('VISUAL', $null, 'User')

$link = Join-Path $HOME "AppData\Local\nvim"
if (Test-Path $link) {
  Remove-Item -Force -Recurse $link
  Write-Host "Removed $link"
}
Write-Host "neovim uninstall done (nvim-config repo left in ~/repos)"
