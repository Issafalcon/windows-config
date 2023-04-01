local wezterm = require("wezterm")
local colors = require("colors")
local fonts = require("fonts")
local windows = require("windows")

local config = {}

if wezterm.config_builder then
    config = wezterm.config_builder()
end

-- Custom Config
colors.apply_to_config(config)
fonts.apply_to_config(config)
windows.apply_to_config(config)

return config