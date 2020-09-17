
param (
    [string]$SpaceCliRoot = "c:\space-cli"
)

$ErrorActionPreference = 'stop'

# Constants
$SpaceCliFileName = "space-cli.exe"
$SpaceCliFilePath = "${SpaceCliRoot}\${SpaceCliFileName}"

# GitHub Org and repo hosting Space CLI
$GitHubOrg="spaceuptech"
$GitHubRepo="space-cloud"

# Set Github request authentication for basic authentication.
if ($Env:GITHUB_USER) {
    $basicAuth = [System.Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes($Env:GITHUB_USER + ":" + $Env:GITHUB_TOKEN));
    $githubHeader = @{"Authorization"="Basic $basicAuth"}
} else {
    $githubHeader = @{}
}

if((Get-ExecutionPolicy) -gt 'RemoteSigned' -or (Get-ExecutionPolicy) -eq 'ByPass') {
    Write-Output "PowerShell requires an execution policy of 'RemoteSigned'."
    Write-Output "To make this change please run:"
    Write-Output "'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    break
}

# Change security protocol to support TLS 1.2 / 1.1 / 1.0 - old powershell uses TLS 1.0 as a default protocol
[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls"

# Check if Space CLI is installed.
if (Test-Path $SpaceCliFilePath -PathType Leaf) {
    Write-Warning "Space cli is detected - $SpaceCliFilePath"
    Invoke-Expression "$SpaceCliFilePath --version"
    Write-Output "Reinstalling Space clid..."
} else {
    Write-Output "Installing Space cli..."
}

# Create Space Directory
Write-Output "Creating $SpaceCliRoot directory"
New-Item -ErrorAction Ignore -Path $SpaceCliRoot -ItemType "directory"
if (!(Test-Path $SpaceCliRoot -PathType Container)) {
    throw "Cannot create $SpaceCliRoot"
}

# Get the list of release from GitHub
$releases = Invoke-RestMethod -Headers $githubHeader -Uri "https://api.github.com/repos/${GitHubOrg}/${GitHubRepo}/releases" -Method Get
if ($releases.Count -eq 0) {
    throw "No releases from github.com/spaceuptech/space-cloud repo"
}

# Filter windows binary and download archive
$windowsAsset = $releases | Where-Object { $_.tag_name } | Select-Object -First 1
if (!$windowsAsset) {
    throw "Cannot find the windows Space CLI binary"
}

$fileName = "space-cli-" + $windowsAsset.tag_name + ".zip"
$zipFilePath = $SpaceCliRoot + "\" + $fileName
Write-Output "Downloading $zipFilePath ..."

# $githubHeader.Accept = "application/octet-stream"
$baseUrl = "https://storage.googleapis.com/space-cloud/windows/" + $fileName
Invoke-WebRequest -Uri $baseUrl -OutFile $zipFilePath
if (!(Test-Path $zipFilePath -PathType Leaf)) {
    throw "Failed to download Space Cli binary - $zipFilePath"
}

# Extract Space CLI to $SpaceCliRoot
Write-Output "Extracting $zipFilePath..."
Microsoft.Powershell.Archive\Expand-Archive -Force -Path $zipFilePath -DestinationPath $SpaceCliRoot
if (!(Test-Path $SpaceCliFilePath -PathType Leaf)) {
    throw "Failed to download Space Cli archieve - $zipFilePath"
}

# Check the Space CLI version
Invoke-Expression "$SpaceCliFilePath --version"

# Clean up zipfile
Write-Output "Clean up $zipFilePath..."
Remove-Item $zipFilePath -Force

# Add SpaceCliRoot directory to User Path environment variable
Write-Output "Trying to add $SpaceCliRoot to User Path Environment variable..."
$UserPathEnvionmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
if($UserPathEnvionmentVar -like '*space-cli*') {
    Write-Output "Path already set, skipping to add $SpaceCliRoot to User Path - $UserPathEnvionmentVar"
} else {
    [System.Environment]::SetEnvironmentVariable("PATH", $UserPathEnvionmentVar + ";$SpaceCliRoot", "User")
    $UserPathEnvionmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
    Write-Output "Added $SpaceCliRoot to User Path - $UserPathEnvionmentVar"
}

Write-Output "`r`nSpace CLI is installed successfully."
Write-Output "To get started with Space, please visit https://learn.spaceuptech.com/space-cloud/basics/ ."
