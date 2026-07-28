param (
  $installationdrive = "C"
)

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$target = Join-Path $scriptDir "themes.gitconfig"
$link = Join-Path $HOME "themes.gitconfig"

Write-Host "Linking $link -> $target"
if (Test-Path $link) { Remove-Item -Force $link }
try {
  New-Item -ItemType HardLink -Path $link -Target $target | Out-Null
  Write-Host "  hardlink created"
} catch {
  Copy-Item -Force $target $link
  Write-Host "  copied (hardlink unavailable): $_"
}

if ($installationdrive -ne "C") {
  Write-Host "Linking C:\Program Files\Git -> ${installationdrive}:\Program Files\Git"
  try {
    New-Item -ItemType SymbolicLink -Path "C:/Program Files/Git/" -Target "${installationdrive}:/Program Files/Git/" -ErrorAction Stop | Out-Null
  } catch {
    Write-Warning "Could not link C:\Program Files\Git (needs admin): $_"
  }
}

Write-Host "Configuring git delta / editor globals..."
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
Write-Host "git config done"
