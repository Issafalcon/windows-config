param (
  $installationdrive = "C", 
  $modulename
)

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

function Install-NeededFor
{
  param(
    [string] $packageName = ''
    , [bool] $defaultAnswer = $true
  )
  if ($packageName -eq '')
  { return $false 
  }
  
  $yes = '6'
  $no = '7'
  $msgBoxTimeout = '-1'
  $defaultAnswerDisplay = 'Yes'
  $buttonType = 0x4;
  if (!$defaultAnswer)
  { $defaultAnswerDisplay = 'No'; $buttonType = 0x104; 
  }
  
  $answer = $msgBoxTimeout
  try
  {
    $timeout = 10
    $question = "Do you need to install $($packageName)? Defaults to `'$defaultAnswerDisplay`' after $timeout seconds"
    $msgBox = New-Object -ComObject WScript.Shell
    $answer = $msgBox.Popup($question, $timeout, "Install $packageName", $buttonType)
  } catch
  {
  }
  
  if ($answer -eq $yes -or ($answer -eq $msgBoxTimeout -and $defaultAnswer -eq $true))
  {
    write-host "Installing $packageName"
    return $true
  }
  
  write-host "Not installing $packageName"
  return $false
}

function InstallModule
{
  param(
    [string] $moduleDir = ''
  )

  if (Test-Path $moduleDir/install.ps1 -PathType Leaf)
  {
    & "${moduleDir}/install.ps1" -installationdrive ${installationDrive}
  }
  if (Test-Path $moduleDir/config.ps1 -PathType Leaf)
  {
    & "${moduleDir}/config.ps1" -installationdrive ${installationDrive}
  }
}

if ($modulename -eq "all")
{
  $moduleOrder = @(
    'git',
    'scoop',
    'wezterm',
    # Neovim module install python3 and node modules as part of it
    'neovim',
    'dotnet',
    'omnisharp',
    'powershell', 
    'lazygit'
  )
  foreach ($module in $moduleOrder)
  {
    if (Install-NeededFor $module)
    {
      InstallModule -moduleDir "${scriptDir}\$module"
    }
  }
} else
{
  InstallModule -moduleDir "${scriptDir}\$modulename"
}
