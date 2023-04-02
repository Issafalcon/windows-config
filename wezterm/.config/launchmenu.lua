local launchmenu = {}

function launchmenu.apply_to_config(config)
    config.launch_menu = {
      {
        label = "Windows - Configurations",
        args = { 
          'wezterm',
          'start'
         }
      }
    }
end

return launchmenu