$weztermConfigDir = Join-Path $HOME "repos\wezterm-config"

if (-not (Test-Path $weztermConfigDir)) {
  Write-Host "Cloning wezterm-config -> $weztermConfigDir"
  git clone https://github.com/Issafalcon/wezterm-config.git $weztermConfigDir
} else {
  Write-Host "wezterm-config already present: $weztermConfigDir"
}

$weztermCustomProjectDir = Join-Path $HOME ".config\wezterm-projects"
if (Test-Path $weztermCustomProjectDir) {
  Write-Host "WezTerm projects folder exists. Skipping creation"
} else {
  Write-Host "Creating $weztermCustomProjectDir"
  New-Item -ItemType Directory -Path $weztermCustomProjectDir | Out-Null
}

$link = Join-Path $HOME ".config\wezterm"
Write-Host "Linking $link -> $weztermConfigDir"
if (Test-Path $link) { Remove-Item -Force -Recurse $link }
New-Item -ItemType Junction -Path $link -Target $weztermConfigDir | Out-Null
Write-Host "wezterm config done"
