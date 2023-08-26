param (
  [string]$installationdrive = "C"
)

$weztermConfigDir = "${installationdrive}:\repos\wezterm-config"

git clone https://github.com/Issafalcon/wezterm-config.git $weztermConfigDir

$weztermCustomProjectDir =  "$HOME/.config/wezterm-projects"
if (Test-Path $weztermCustomProjectDir)
{
  "WezTerm projects folder exists. Skipping creation"
} else
{
  New-Item -ItemType Directory -Path $weztermCustomProjectDir
}

New-Item -ItemType SymbolicLink -Path "$HOME/.config/wezterm/" -Target $weztermConfigDir/
