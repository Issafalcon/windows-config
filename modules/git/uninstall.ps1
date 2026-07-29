# Git itself is a TUI prerequisite — only remove delta + config this module added.
scoop uninstall delta

$link = Join-Path $HOME "themes.gitconfig"
if (Test-Path $link) {
  Remove-Item -Force $link
  Write-Host "Removed $link"
}

git config --global --unset core.pager 2>$null
git config --global --unset include.path 2>$null
git config --global --unset interactive.diffFilter 2>$null
git config --global --unset-all delta.navigate 2>$null
git config --global --unset-all delta.line-numbers 2>$null
git config --global --unset-all delta.side-by-side 2>$null
git config --global --unset-all delta.syntax-theme 2>$null
git config --global --unset-all delta.features 2>$null
Write-Host "git delta config cleared"
