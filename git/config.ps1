param (
  $installationdrive = "C"
)

if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments -Wait
  Break
}

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

# Symlink themes config to correct location
New-Item -ItemType SymbolicLink -Path "~/themes.gitconfig" -Target "${scriptDir}/themes.gitconfig" -Force

if ($installationdrive -ne "C")
{
  New-Item -ItemType SymbolicLink -Path "C:/Program Files/Git/" -Target "${installationdrive}:/Program Files/Git/"
}

# Set delta config
git config --global core.pager "delta --dark --paging=never"
git config --global include.path "~/themes.gitconfig"
git config --global interactive.diffFilter "delta --color-only"
git config --global delta.navigate "true"
git config --global delta.line-numbers "true"
git config --global delta.side-by-side "false"
git config --global delta.syntax-theme "Dracula"
git config --global delta.features "decorations line-numbers zebra-dark"
git config --global merge.conflictstyle "diff3"
git config --global core.editor nvim

Read-Host -Prompt "Press Enter to exit"
