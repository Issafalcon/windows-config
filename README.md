# Windows Config

Collection of configuration files and scripts to set up my dev environment on Windows 10+ machines

## Prerequisites

1. Download and install git for windows (onto your drive of choice)
2. Create `~/repos` and clone this repo into it as `windows-config`
3. Run the `module-install.ps1` script to install modules (see script for specific order requried)

## Usage

For now, this repo needs to be cloned into the `~/repos` folder.

Modules can be installed and configured using the base `./module-install.ps1` script.

To install all the modules, run the script with the `-modulename` parameters set to `all`.

```
.\module-install.ps1 -installationdrive D -modulename all
```

You will be prompted which modules you wish to install, and then be able to check the output of each install and config script before continuing (some will run in elevated powershell prompt).

Individual modules can be installed using the specific module name (named after folder name).
