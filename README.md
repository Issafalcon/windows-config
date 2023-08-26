# Windows Config

Collection of configuration files and scripts to set up my dev environment on Windows 10+ machines

## Usage

### Setup

First, run `dev-setup.ps1` to create a repos folder and symbolic links to `$HOME` i.e.

```
Start-Process powershell -verb runas -ArgumentList "-file .\dev-setup.ps1 -installationdrive D -repositoriespath repos"
```

This will create the `D:\repos` path, and symlink it to `~/repos` and install scoop

### Install scoop
Modules can be installed and configured using the base `./module-install.ps1` script.

