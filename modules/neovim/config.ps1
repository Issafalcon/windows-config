# Make Neovim the default editor (User scope — no admin)
[System.Environment]::SetEnvironmentVariable('EDITOR', "nvim", "User")
[System.Environment]::SetEnvironmentVariable('VISUAL', "nvim", "User")

$nvimConfigDir = "${HOME}\repos\nvim-config"

if (-not (Test-Path $nvimConfigDir)) {
  git clone https://github.com/Issafalcon/nvim-config.git $nvimConfigDir
}

New-Item -ItemType SymbolicLink -Path "~/AppData/Local/nvim/" -Target $nvimConfigDir/ -Force
