# Make Neovim the default editor (User scope — no admin)
Write-Host "Setting EDITOR/VISUAL=nvim (User)"
[System.Environment]::SetEnvironmentVariable('EDITOR', "nvim", "User")
[System.Environment]::SetEnvironmentVariable('VISUAL', "nvim", "User")

$nvimConfigDir = Join-Path $HOME "repos\nvim-config"

if (-not (Test-Path $nvimConfigDir)) {
  Write-Host "Cloning nvim-config -> $nvimConfigDir"
  git clone https://github.com/Issafalcon/nvim-config.git $nvimConfigDir
} else {
  Write-Host "nvim-config already present: $nvimConfigDir"
}

$link = Join-Path $HOME "AppData\Local\nvim"
Write-Host "Linking $link -> $nvimConfigDir"
if (Test-Path $link) { Remove-Item -Force -Recurse $link }
New-Item -ItemType Junction -Path $link -Target $nvimConfigDir | Out-Null
Write-Host "neovim config done"
