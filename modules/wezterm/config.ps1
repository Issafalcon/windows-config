$weztermConfigDir = Join-Path $HOME "repos\wezterm-config"

if (-not (Test-Path $weztermConfigDir)) {
  git clone https://github.com/Issafalcon/wezterm-config.git $weztermConfigDir
}

$weztermCustomProjectDir = Join-Path $HOME ".config\wezterm-projects"
if (Test-Path $weztermCustomProjectDir) {
  "WezTerm projects folder exists. Skipping creation"
} else {
  New-Item -ItemType Directory -Path $weztermCustomProjectDir | Out-Null
}

$link = Join-Path $HOME ".config\wezterm"
if (Test-Path $link) { Remove-Item -Force -Recurse $link }
# Directory junction: no admin / Developer Mode required.
New-Item -ItemType Junction -Path $link -Target $weztermConfigDir | Out-Null
