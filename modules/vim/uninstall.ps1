foreach ($name in @('.ideavimrc', '.vimrc')) {
  $link = Join-Path $HOME $name
  if (Test-Path $link) {
    Remove-Item -Force $link
    Write-Host "Removed $link"
  }
}
