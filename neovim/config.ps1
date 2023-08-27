if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments -Wait
  Break
}

# Make Neovim the default editor
[System.Environment]::SetEnvironmentVariable('EDITOR', "nvim", "Machine")
[System.Environment]::SetEnvironmentVariable('VISUAL', "nvim", "Machine")

$nvimConfigDir = "${HOME}\repos\nvim-config"

git clone https://github.com/Issafalcon/nvim-config.git $nvimConfigDir

New-Item -ItemType SymbolicLink -Path "~/AppData/Local/nvim/" -Target $nvimConfigDir/

Read-Host -Prompt "Press Enter to exit"
