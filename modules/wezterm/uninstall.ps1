scoop uninstall wezterm

$link = Join-Path $HOME ".config\wezterm"
if (Test-Path $link) {
  Remove-Item -Force -Recurse $link
  Write-Host "Removed $link"
}
Write-Host "wezterm-config repo left in ~/repos"
