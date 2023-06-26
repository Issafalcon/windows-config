# Windows Config

Collection of configuration files and scripts to set up my dev environment on Windows 10+ machines

## Usage

Modules can be installed and configured using the base `./module-install.ps1` script.

## Terminal

Contains the Windows Terminal software settings in `settings.json`

- Copy settings to override the file in `%USERPROFILE%\AppData\Local\Packages\<Windows terminal package location>`

## Lazygit

Contains config file for Windows

- Copy this file to %APPDATA%\Roaming\lazygit

## Omnisharp

Contains the configuration for the omnisharp-roslyn server

- Copy the `omnisharp.json` file to `%USERPROFILE%/.omnisharp` folder
