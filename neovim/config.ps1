param (
  $installationdrive = "C"
)

$nvimConfigDir = "${installationdrive}:/repos/nvim-config"

mkdir ${installationdrive}:/repos/
git clone https://github.com/Issafalcon/nvim-config.git $nvimConfigDir

New-Item -ItemType SymbolicLink -Path "~/AppData/Local/nvim/" -Target $nvimConfigDir/
