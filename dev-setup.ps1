###################################################################
# Prerequisites for environment setup
###################################################################

param (
  $installationdrive = "C",
  $repositoriespath = "repos"
)

# Create git repository folder (if doesn't exist) and symlink to home if using a different install drive
if (Test-Path -Path "${installationdrive}:/${repositoriespath}")
{
  "${installationdrive}:/${repositoriespath} folder exists. Skipping creation"
} else
{
  New-Item -ItemType Directory -Path "${installationdrive}:/${repositoriespath}"
}
