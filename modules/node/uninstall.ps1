# nvm-windows was installed via nvm-setup.exe (not winget). Best-effort cleanup:
$nvmRoot = [Environment]::GetEnvironmentVariable('NVM_HOME', 'User')
if (-not $nvmRoot) { $nvmRoot = Join-Path ${env:ProgramFiles} 'nvm' }

$unins = @(
  (Join-Path $nvmRoot 'unins000.exe'),
  (Join-Path ${env:ProgramFiles} 'nvm\unins000.exe')
) | Where-Object { Test-Path $_ } | Select-Object -First 1

if ($unins) {
  Write-Host "Running $unins"
  Start-Process -FilePath $unins -ArgumentList '/VERYSILENT' -Wait
} else {
  Write-Warning "nvm uninstaller not found; remove 'NVM for Windows' from Apps & Features if needed."
}
