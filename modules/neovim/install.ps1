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

# nvim-treesitter (main) shells out to `tree-sitter build` → cc-rs, which on
# Windows defaults to cl.exe. VS may be installed without cl on PATH; scoop
# mingw's gcc is what we ship, so point the compiler at it.
[System.Environment]::SetEnvironmentVariable('CC', 'gcc', 'User')
[System.Environment]::SetEnvironmentVariable('CXX', 'g++', 'User')

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

# win32yank for clipboard support
scoop install win32yank

# python3 / node come from module.yaml dependencies (install via TUI queue)

# npm v12+ blocks dependency install scripts unless allowlisted (global installs
# use --allow-scripts; tree-sitter-cli needs its install.js to fetch the binary).
npm install -g tree-sitter-cli --allow-scripts=tree-sitter-cli
