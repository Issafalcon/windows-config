param([Parameter(Mandatory)][string]$JobDir)

$logRoot = Join-Path $env:LOCALAPPDATA "windows-config-tui"
New-Item -ItemType Directory -Force -Path $logRoot | Out-Null
$bootLog = Join-Path $logRoot "elevated-helper.log"

function Write-BootLog([string]$Message) {
    $line = "{0:o} {1}" -f (Get-Date).ToUniversalTime(), $Message
    Add-Content -Path $bootLog -Value $line
    Write-Host $line
}

trap {
    Write-BootLog "FATAL: $_"
    Write-BootLog $_.ScriptStackTrace
    Start-Sleep -Seconds 20
    break
}

try {
    Write-BootLog "starting JobDir=$JobDir pid=$PID"
    New-Item -ItemType Directory -Force -Path $JobDir | Out-Null
    Set-Content -NoNewline -Path (Join-Path $JobDir "ready") -Value "1"
    Write-BootLog "ready"
} catch {
    Write-BootLog "startup failed: $_"
    Start-Sleep -Seconds 20
    throw
}

# Elevated RunAs sessions often miss User PATH (where Scoop/Git shims live).
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
        # Do not use $args — it is automatic/read-only inside scriptblocks.
        $scriptArgs = @()
        if ($null -ne $request.args) { $scriptArgs = @($request.args) }
        Write-BootLog "job $($request.id) -> $($request.scriptPath)"
        try {
            # Nested pwsh so script `exit` cannot kill this helper session.
            & pwsh -NoProfile -ExecutionPolicy Bypass -File $request.scriptPath @scriptArgs *>&1 |
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
        Write-BootLog "job $($request.id) exit=$exitCode"
        Set-Content -NoNewline -Path (Join-Path $JobDir "$($request.id).done") -Value $exitCode
        Remove-Item -Force $_.FullName
    }
    Start-Sleep -Milliseconds 100
}

Write-BootLog "shutdown"
