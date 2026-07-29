# Official Scoop uninstall (removes scoop root + shims).
$uninstall = Join-Path $env:TEMP 'scoop-uninstall.ps1'
Invoke-RestMethod 'https://raw.githubusercontent.com/ScoopInstaller/Install/master/uninstall.ps1' -OutFile $uninstall
& $uninstall -RunAsAdmin:$false
Remove-Item $uninstall -Force -ErrorAction SilentlyContinue
