# Windows Config

Collection of configuration files and scripts for a Windows 10+ dev environment, managed by a native **Windows Config TUI**.

## Prerequisites

1. **Git for Windows** (install manually first)
2. **PowerShell 7** (`pwsh`) — `winget install Microsoft.PowerShell`
3. Clone this repo (typically `~/repos/windows-config`)
4. Optional but recommended: enable **Windows Developer Mode** (Settings → Privacy & security → For developers) so user-profile **symbolic links** work without elevation

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

- `admin_scripts: [config.ps1]` — elevate that script only (typical for **symbolic links**); `install.ps1` stays user-scoped (scoop-friendly)
- `requires_admin: true` — elevate both `install.ps1` and `config.ps1` when `admin_scripts` is empty

Elevated work goes through a session helper (`tools/ElevatedHelper.ps1`):

- First admin job → one UAC prompt (`pwsh -Verb RunAs`)
- Later admin jobs in the same TUI session reuse the helper (no more UAC)
- Quitting the TUI shuts the helper down

Without elevation, user-profile symlinks need **Windows Developer Mode**. Modules that create links list `admin_scripts: [config.ps1]` so the TUI elevates those steps for you.

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
