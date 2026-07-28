param([Parameter(Mandatory)][string]$JobDir)

New-Item -ItemType Directory -Force -Path $JobDir | Out-Null
Set-Content -NoNewline -Path (Join-Path $JobDir "ready") -Value "1"

while (-not (Test-Path (Join-Path $JobDir "shutdown"))) {
    Get-ChildItem -Path $JobDir -Filter "*.req.json" -ErrorAction SilentlyContinue | ForEach-Object {
        $request = Get-Content -Raw $_.FullName | ConvertFrom-Json
        $logPath = Join-Path $JobDir "$($request.id).log"
        $args = @()
        if ($null -ne $request.args) { $args = @($request.args) }
        try {
            & pwsh -NoProfile -File $request.scriptPath @args *>&1 |
                ForEach-Object { $_.ToString() | Add-Content -Path $logPath }
            $exitCode = $LASTEXITCODE
            if ($null -eq $exitCode) { $exitCode = 0 }
        } catch {
            $_.ToString() | Add-Content -Path $logPath
            $exitCode = 1
        }
        Set-Content -NoNewline -Path (Join-Path $JobDir "$($request.id).done") -Value $exitCode
        Remove-Item -Force $_.FullName
    }
    Start-Sleep -Milliseconds 100
}
