param (
  $installationdrive = "C"
)

# Neovim
scoop bucket add extras
scoop install vcredist2022
scoop install neovim

if ($installationdrive -ne "C") {
  # If installation drive is specified, move the heavy data files into this drive instead
  [System.Environment]::SetEnvironmentVariable('XDG_DATA_HOME', "${installationdrive}:\AppData", "User")
}

# Install gcc (also adds 7zip)
scoop bucket add main
scoop install mingw

# Go
scoop install go

## RipGrep, fd, fzf and Silver Searcher for finding and sorting
scoop install ripgrep
scoop install ag
scoop install fzf
scoop install fd

# Tools for mason package manager extraction
scoop install wget
scoop install unzip
scoop install gzip
