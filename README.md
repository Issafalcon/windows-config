# Windows Config

Collection of configuration files and scripts for a Windows 10+ dev environment, managed by a native **Windows Config TUI**.

## Prerequisites

1. **Git for Windows** (install manually first)
2. **PowerShell 7** (`pwsh`) — `winget install Microsoft.PowerShell`
3. Clone this repo (typically `~/repos/windows-config`)
4. Optional: enable **Windows Developer Mode** only if you still want classic symbolic links; the bundled modules use hard links / junctions instead

## TUI (preferred)

### Prebuilt exe (no Go)

On each `v*` tag, CI publishes a zip under [Releases](https://github.com/Issafalcon/windows-config/releases):

1. Download `windows-config-tui-*-windows-amd64.zip` (or `arm64`)
2. Unzip (keep `tools\ElevatedHelper.ps1` beside the exe)
3. Clone this repo (for `modules/`) or set `WINDOWS_CONFIG_MODULES_DIR`
4. Run `windows-config-tui.exe` from Windows Terminal / PowerShell 7

### Build from source

```powershell
go build -o windows-config-tui.exe ./tui/
.\windows-config-tui.exe
```

Or `make -C tui windows`.

The TUI discovers packages under `modules/*/module.yaml`, plans dependency order, runs `install.ps1` then `config.ps1`, and tracks installs in `%USERPROFILE%\.windowsConfigModules`.

Override the modules path with `WINDOWS_CONFIG_MODULES_DIR` or the first-run path prompt (`~/.config/windows-config-tui/config.yaml`).

### Elevation model

The TUI always starts **unelevated**. Privilege is requested only for scripts that need it:

- `admin_scripts: [config.ps1]` — elevate that script only (rare machine-wide changes)
- `requires_admin: true` — elevate both `install.ps1` and `config.ps1` when `admin_scripts` is empty

Most profile config uses **hard links** (files) or **junctions** (directories), which do **not** need admin or Developer Mode. Scoop-based `install.ps1` scripts always stay unelevated — Scoop refuses / breaks under admin.

When elevation is required, work goes through a session helper (`tools/ElevatedHelper.ps1`):

- First admin job → one UAC prompt (`pwsh -Verb RunAs`, ExecutionPolicy Bypass)
- Later admin jobs in the same TUI session reuse the helper (no more UAC)
- Quitting the TUI shuts the helper down
- Startup/job traces: `%LOCALAPPDATA%\windows-config-tui\elevated-helper.log`

### Layout

```
windows-config/
  tui/                 # Go TUI entrypoint
  internal/            # TUI packages
  tools/
    ElevatedHelper.ps1 # session elevated job runner
  modules/
    <name>/
      module.yaml
      install.ps1      # optional
      config.ps1       # optional (symlinks / profile — no Stow)
  module-install.ps1   # deprecated CLI fallback
```

## Deprecated: `module-install.ps1`

`.\module-install.ps1` still works and now looks under `modules\`, but the TUI is the supported installer. Example:

```powershell
.\module-install.ps1 -installationdrive D -modulename all
.\module-install.ps1 -modulename neovim
```

Scripts no longer self-elevate with `Start-Process -Verb RunAs`; elevation is the TUI helper’s job via `admin_scripts` / `requires_admin`.
