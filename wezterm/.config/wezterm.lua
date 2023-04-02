local wezterm = require("wezterm")
local colors = require("colors")
local fonts = require("fonts")
local windows = require("windows")
local launchmenu = require("launchmenu")
local keybindings = require("keybindings")
require("eventhandlers")

local config = {}

if wezterm.config_builder then
    config = wezterm.config_builder()
end

-- Custom Config
config.default_prog = { 'pwsh' }
config.use_fancy_tab_bar = true
config.disable_default_key_bindings = true
config.enable_csi_u_key_encoding = true
config.leader = { key = "Space", mods = "CTRL|SHIFT" }
config.keys = keybindings.create_keybinds()
config.selection_word_boundary = " \t\n{}[]()\"'`,;:│=&!%"

colors.apply_to_config(config)
fonts.apply_to_config(config)
windows.apply_to_config(config)
launchmenu.apply_to_config(config)

return config