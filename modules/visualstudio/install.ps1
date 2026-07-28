# Visual Studio Enterprise via winget, with ASP.NET / Web API workload + Copilot.
# Extensions (Reqnroll, AWS Toolkit) are installed as VSIX after setup finishes.
#
# Docs: https://learn.microsoft.com/en-us/visualstudio/install/use-command-line-parameters-to-install-visual-studio

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$vsWingetId = 'Microsoft.VisualStudio.Enterprise'
# NetWeb = ASP.NET and web development (ASP.NET Core Web API templates, IIS Express, etc.)
# includeRecommended pulls Copilot-adjacent web tooling; Copilot is also added explicitly.
$vsOverride = @(
  '--wait'
  '--passive'
  '--addProductLang'
  'En-us'
  '--add'
  'Microsoft.VisualStudio.Workload.NetWeb;includeRecommended'
  '--add'
  'Component.VisualStudio.GitHub.Copilot'
) -join ' '

function Get-VsWhere {
  $p = Join-Path ${env:ProgramFiles(x86)} 'Microsoft Visual Studio\Installer\vswhere.exe'
  if (-not (Test-Path $p)) { throw "vswhere.exe not found at $p (Visual Studio Installer missing)" }
  return $p
}

function Get-EnterpriseInstallPath {
  $vswhere = Get-VsWhere
  $path = & $vswhere -products Microsoft.VisualStudio.Product.Enterprise -latest -property installationPath 2>$null
  if ($path) { return ($path | Select-Object -First 1) }
  return $null
}

function Install-OrModifyVisualStudio {
  Write-Host "Visual Studio Enterprise via winget ($vsWingetId)"
  Write-Host "  workloads: NetWeb;includeRecommended + GitHub Copilot"
  Write-Host "  override: $vsOverride"

  $wingetArgs = @(
    'install'
    '--id', $vsWingetId
    '-e'
    '--source', 'winget'
    '--accept-package-agreements'
    '--accept-source-agreements'
    '--disable-interactivity'
    '--override', $vsOverride
  )
  & winget @wingetArgs
  $code = $LASTEXITCODE
  if ($code -eq 0) { return }

  # Already installed / winget refused a fresh install — modify the existing instance.
  Write-Host "winget install exited $code; modifying existing Enterprise install if present..."
  $installPath = Get-EnterpriseInstallPath
  if (-not $installPath) {
    throw "Visual Studio Enterprise is not installed and winget install failed (exit $code)"
  }
  $setup = Join-Path ${env:ProgramFiles(x86)} 'Microsoft Visual Studio\Installer\setup.exe'
  Write-Host "Modifying $installPath"
  $modifyArgs = @(
    'modify'
    '--installPath', $installPath
    '--wait'
    '--passive'
    '--add', 'Microsoft.VisualStudio.Workload.NetWeb;includeRecommended'
    '--add', 'Component.VisualStudio.GitHub.Copilot'
  )
  $p = Start-Process -FilePath $setup -ArgumentList $modifyArgs -Wait -PassThru
  if ($p.ExitCode -ne 0) {
    throw "Visual Studio modify failed (exit $($p.ExitCode))"
  }
}

function Get-MarketplaceVsixUrl([string]$Publisher, [string]$Extension) {
  $body = @{
    filters = @(
      @{
        criteria = @(
          @{ filterType = 7; value = "$Publisher.$Extension" }
        )
        pageNumber = 1
        pageSize   = 1
      }
    )
    # IncludeFiles | IncludeVersionProperties | IncludeMetadata ...
    flags = 0x200 + 0x100 + 0x80 + 0x20 + 0x10 + 0x2
  } | ConvertTo-Json -Depth 6

  $resp = Invoke-RestMethod `
    -Method Post `
    -Uri 'https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery?api-version=7.2-preview.1' `
    -ContentType 'application/json' `
    -Body $body

  $ext = $resp.results[0].extensions[0]
  if (-not $ext) { throw "Marketplace extension not found: $Publisher.$Extension" }
  $ver = $ext.versions[0]
  $file = $ver.files | Where-Object {
    $_.assetType -eq 'Microsoft.VisualStudio.Services.VSIXPackage' -or
    $_.assetType -like '*.vsix'
  } | Select-Object -First 1
  if (-not $file -or -not $file.source) {
    throw "No VSIX asset for $Publisher.$Extension@$($ver.version)"
  }
  Write-Host "  $($ext.displayName) $($ver.version)"
  return $file.source
}

function Install-VsixFromMarketplace([string]$Publisher, [string]$Extension) {
  $vswhere = Get-VsWhere
  $vsixInstaller = & $vswhere -latest -products * -find 'Common7\IDE\VSIXInstaller.exe' | Select-Object -First 1
  if (-not $vsixInstaller) {
    throw 'VSIXInstaller.exe not found (is Visual Studio installed?)'
  }

  Write-Host "Downloading $Publisher.$Extension ..."
  $url = Get-MarketplaceVsixUrl $Publisher $Extension
  $vsix = Join-Path $env:TEMP "$Publisher.$Extension.vsix"
  Invoke-WebRequest -Uri $url -OutFile $vsix -UseBasicParsing

  Write-Host "Installing $vsix via VSIXInstaller..."
  # /admin = all users; /quiet = no UI; shut down is avoided with /skuName:Enterprise when possible
  $p = Start-Process -FilePath $vsixInstaller -ArgumentList @('/quiet', '/admin', $vsix) -Wait -PassThru
  Remove-Item -Force $vsix -ErrorAction SilentlyContinue
  if ($p.ExitCode -ne 0 -and $p.ExitCode -ne 1001) {
    # 1001 = already installed
    throw "VSIXInstaller failed for $Publisher.$Extension (exit $($p.ExitCode))"
  }
  if ($p.ExitCode -eq 1001) {
    Write-Host "  already installed"
  } else {
    Write-Host "  installed"
  }
}

Install-OrModifyVisualStudio

$installPath = Get-EnterpriseInstallPath
if (-not $installPath) {
  throw 'Visual Studio Enterprise install path not found after setup'
}
Write-Host "Installed at: $installPath"

Write-Host 'Installing marketplace extensions...'
Install-VsixFromMarketplace 'Reqnroll' 'ReqnrollForVisualStudio2022'
Install-VsixFromMarketplace 'AmazonWebServices' 'AWSToolkitforVisualStudio2022'

Write-Host 'Visual Studio Enterprise setup complete.'
