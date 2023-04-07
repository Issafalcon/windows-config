param (
  $installationdrive = "C"
)

$weztermConfigDir = "${installationdrive}:/repos/wezterm-config"
$weztermCustomProjectDir =  "~/.config/wezterm-projects"

mkdir ${installationdrive}:/repos/
git clone https://github.com/Issafalcon/wezterm-config.git $weztermConfigDir

if (Test-Path $weztermCustomProjectDir)
{
  "WezTerm projects folder exists. Skipping creation"
} else
{
  New-Item -ItemType Directory -Path $weztermCustomProjectDir
}

New-Item -ItemType SymbolicLink -Path "~/.config/wezterm" -Target $weztermConfigDir
