scoop uninstall lazygit

$link = Join-Path $HOME "AppData\Roaming\lazygit\config.yml"
if (Test-Path $link) {
  Remove-Item -Force $link
  Write-Host "Removed $link"
}
