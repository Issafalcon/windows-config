$link = Join-Path $HOME ".omnisharp\omnisharp.json"
if (Test-Path $link) {
  Remove-Item -Force $link
  Write-Host "Removed $link"
}
