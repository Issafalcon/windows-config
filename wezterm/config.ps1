param (
  $installationdrive = "C"
)

$weztermConfigDir = "${installationdrive}:/repos/wezterm-config"

mkdir ${installationdrive}:/repos/
git clone https://github.com/Issafalcon/wezterm-config.git $weztermConfigDir

New-Item -ItemType SymbolicLink -Path "~/.config/wezterm" -Target $weztermConfigDir
