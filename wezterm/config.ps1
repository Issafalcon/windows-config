param (
  $installationdrive = "default"
)

if ($installationdrive -eq "default")
{
  $weztermConfigDir = "~/repos/wezterm"
} else
{
  $weztermConfigDir = "${installationdrive}:/repos/wezterm"
}
$weztermCustomProjectDir =  "~/.config/wezterm-projects"

git clone https://github.com/Issafalcon/wezterm-config.git $weztermConfigDir

if (Test-Path $weztermCustomProjectDir)
{
  "WezTerm projects folder exists. Skipping creation"
} else
{
  New-Item -ItemType Directory -Path $weztermCustomProjectDir
}

New-Item -ItemType SymbolicLink -Path "~/.config/wezterm" -Target $weztermConfigDir
