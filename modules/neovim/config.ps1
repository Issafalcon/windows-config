# Make Neovim the default editor (User scope — no admin)
[System.Environment]::SetEnvironmentVariable('EDITOR', "nvim", "User")
[System.Environment]::SetEnvironmentVariable('VISUAL', "nvim", "User")

$nvimConfigDir = Join-Path $HOME "repos\nvim-config"

if (-not (Test-Path $nvimConfigDir)) {
  git clone https://github.com/Issafalcon/nvim-config.git $nvimConfigDir
}

$link = Join-Path $HOME "AppData\Local\nvim"
if (Test-Path $link) { Remove-Item -Force -Recurse $link }
# Directory junction: no admin / Developer Mode required.
New-Item -ItemType Junction -Path $link -Target $nvimConfigDir | Out-Null
