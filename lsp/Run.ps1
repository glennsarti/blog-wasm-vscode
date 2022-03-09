param([Switch]$Windows, [Switch]$Wasm)

# Set Default
if (!($Windows.IsPresent -or $Wasm.IsPresent)) { $Windows = $true }

if ($Windows.IsPresent) {
  .\Build.ps1 -Windows

  Write-Host "Running..."
  .\pkg\demo-ls.exe serve
} elseif ($WASM.IsPresent) {
  .\Build.ps1 -Wasm

  if ($null -eq $Global:GOROOT) {
    $Global:GOROOT = & go env GOROOT
  }
  $wasmexec = join-Path $Global:GOROOT "/misc/wasm/wasm_exec.js"

  Write-Host "Running..."
  & node $wasmexec pkg/demo-ls.wasm
}
