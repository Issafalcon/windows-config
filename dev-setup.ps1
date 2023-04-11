###################################################################
# Prerequisites for environment setup
###################################################################

param (
  $installationdrive = "default",
  $repositoriespath = "repos"
)

# Create git repository folder (if doesn't exist) and symlink to home if using a different install drive
if ($installationdrive -eq "default")
{
  if (Test-Path -Path "~/${repositoriespath}")
  {
    "~/${repositoriespath} folder exists. Skipping creation"
  } else
  {
    New-Item -ItemType Directory -Path "~/${repositoriespath}"
  }
} else
{
  if (Test-Path -Path "${installationdrive}:/${repositoriespath}")
  {
    "${installationdrive}:/${repositoriespath} folder exists. Skipping creation"
  } else
  {
    New-Item -ItemType Directory -Path "${installationdrive}:/${repositoriespath}"
  }

  New-Item -ItemType SymbolicLink -Path "~/${repositoriespath}" -Target "${installationdrive}:/${repositoriespath}" -Force
}
