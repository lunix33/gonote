##
# This is for Windows ONLY.
# Be sure to enable your Powershell scripts.
# This script set the environment variable needed to work on this project.
#
# GOPATH : (The current path.)
# CGO_ENABLED : 1
##

[System.Environment]::SetEnvironmentVariable(
    'GOPATH', $PSScriptRoot,
    [System.EnvironmentVariableTarget]::User);
[System.Environment]::SetEnvironmentVariable(
    'CGO_ENABLED', 1,
    [System.EnvironmentVariableTarget]::User);