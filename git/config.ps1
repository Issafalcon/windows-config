param (
  [Parameter(Mandatory)]$installationdrive
)

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

# Symlink themes config to correct location
New-Item -ItemType SymbolicLink -Path "~/themes.gitconfig" -Target "${scriptDir}/themes.gitconfig"

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