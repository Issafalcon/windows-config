param (
  $installationdrive = "C"
)

# Git
winget install --id Git.Git -e --source winget --location "${installationdrive}:\Program Files\Git"
[System.Environment]::SetEnvironmentVariable('Path', $env:Path + ";${installationdrive}:\Program Files\Git\cmd", "Machine")
scoop install delta
