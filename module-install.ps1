param ([Parameter(Mandatory)]$installationdrive, $modulename)

$scriptDir = $(Split-Path -parent $MyInvocation.MyCommand.Definition)

function Install-NeededFor {
  param(
    [string] $packageName = ''
    , [bool] $defaultAnswer = $true
  )
  if ($packageName -eq '') { return $false }
  
  $yes = '6'
  $no = '7'
  $msgBoxTimeout = '-1'
  $defaultAnswerDisplay = 'Yes'
  $buttonType = 0x4;
  if (!$defaultAnswer) { $defaultAnswerDisplay = 'No'; $buttonType = 0x104; }
  
  $answer = $msgBoxTimeout
  try {
    $timeout = 10
    $question = "Do you need to install $($packageName)? Defaults to `'$defaultAnswerDisplay`' after $timeout seconds"
    $msgBox = New-Object -ComObject WScript.Shell
    $answer = $msgBox.Popup($question, $timeout, "Install $packageName", $buttonType)
  }
  catch {
  }
  
  if ($answer -eq $yes -or ($answer -eq $msgBoxTimeout -and $defaultAnswer -eq $true)) {
    write-host "Installing $packageName"
    return $true
  }
  
  write-host "Not installing $packageName"
  return $false
}

if ($modulename -eq "all") {
  Get-ChildItem -Recurse -Directory | ForEach-Object {
    if (Install-NeededFor $_.Name) {
      & "${_.FullName}/install.ps1" -installationdrive ${installationDrive}
      & "${_.FullName}/config.ps1"
    }
  }
}
else {
  $modulePath = "${scriptDir}\$modulename"
  & "${modulePath}\install.ps1" -installationdrive ${installationDrive}
  & "${modulePath}\config.ps1" -installationdrive ${installationDrive}
}