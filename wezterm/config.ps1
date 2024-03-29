if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments -Wait
  Break
}

$weztermConfigDir = "${HOME}\repos\wezterm-config"

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

Read-Host -Prompt "Press Enter to exit"
