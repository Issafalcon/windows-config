###################################################################
# Prerequisites for environment setup
###################################################################

param (
  $installationdrive = "C",
  $repositoriespath = "repos"
)

if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator"))  
{  
  $arguments = "& '" +$myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments
  Break
}

# Create git repository folder (if doesn't exist) and symlink to home if using a different install drive
if (Test-Path -Path "${installationdrive}:/${repositoriespath}")
{
  "${installationdrive}:/${repositoriespath} folder exists. Skipping creation"
} else
{
  New-Item -ItemType Directory -Path "${installationdrive}:/${repositoriespath}"
}

# Create the symlinked folder
if (Test-Path -Path "~/${repositoriespath}")
{
  "~/${repositoriespath} folder exists. Skipping creation"
} else
{
  New-Item -ItemType Directory -Path "~/${repositoriespath}"
}
