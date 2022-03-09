param([Switch]$Windows, [Switch]$Wasm)

if (-not(Test-Path 'pkg')) { New-Item -Path 'pkg' -ItemType 'Directory' }

# Set Default
if (!($Windows.IsPresent -or $Wasm.IsPresent)) { $Windows = $true }

# go build
if ($Windows.IsPresent) {
  $ENV:GOOS = 'windows'
  $ENV:GOARCH = 'amd64'
  Write-Host "Building for Windows..."
  & go build -o pkg/demo-ls.exe
}

if ($Wasm.IsPresent) {
  $ENV:GOOS = 'js'
  $ENV:GOARCH = 'wasm'
  Write-Host "Building for Wasm..."
  & go build -o pkg/demo-ls.wasm
}
