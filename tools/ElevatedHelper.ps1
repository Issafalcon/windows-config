param([Parameter(Mandatory)][string]$JobDir)

try {
    New-Item -ItemType Directory -Force -Path $JobDir | Out-Null
    Set-Content -NoNewline -Path (Join-Path $JobDir "ready") -Value "1"
} catch {
    Write-Host "ElevatedHelper failed to start: $_"
    Write-Host "JobDir=$JobDir"
    Start-Sleep -Seconds 15
    throw
}

# Elevated RunAs sessions often miss User PATH (where Scoop shims live).
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath) {
    $env:Path = "$userPath;$env:Path"
}
$scoop = [Environment]::GetEnvironmentVariable('SCOOP', 'User')
if ($scoop) { $env:SCOOP = $scoop }

while (-not (Test-Path (Join-Path $JobDir "shutdown"))) {
    Get-ChildItem -Path $JobDir -Filter "*.req.json" -ErrorAction SilentlyContinue | ForEach-Object {
        $request = Get-Content -Raw $_.FullName | ConvertFrom-Json
        $logPath = Join-Path $JobDir "$($request.id).log"
        $args = @()
        if ($null -ne $request.args) { $args = @($request.args) }
        try {
            # Nested pwsh so script `exit` cannot kill this helper session.
            & pwsh -NoProfile -File $request.scriptPath @args *>&1 |
                ForEach-Object {
                    $line = $_.ToString()
                    Add-Content -Path $logPath -Value $line
                    Write-Host $line
                }
            $exitCode = $LASTEXITCODE
            if ($null -eq $exitCode) { $exitCode = 0 }
        } catch {
            $_.ToString() | Tee-Object -FilePath $logPath -Append | Write-Host
            $exitCode = 1
        }
        Set-Content -NoNewline -Path (Join-Path $JobDir "$($request.id).done") -Value $exitCode
        Remove-Item -Force $_.FullName
    }
    Start-Sleep -Milliseconds 100
}
