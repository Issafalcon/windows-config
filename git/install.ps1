param (
  $installationdrive = "C"
)

# Git - Install isn't needed as git needs to be installed manually first anyway as a pre-requisite
# winget install --id Git.Git -e --source winget --location "${installationdrive}:\Program Files\Git"
# [System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\Program Files\Git\cmd", "Machine")
scoop install delta
