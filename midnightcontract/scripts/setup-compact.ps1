param(
  [string]$CompactPath = "/home/ekko/.local/bin/compact"
)

function Invoke-WSL($ArgsArray) {
  & wsl @ArgsArray 2>&1 | Out-String
}

Write-Host "Checking Compact version..."
$ver = Invoke-WSL @($CompactPath, '--version')
Write-Host $ver

Write-Host "Listing available compilers..."
$list = Invoke-WSL @($CompactPath, 'list', 'compilers')
Write-Host $list

$compilerName = ($list -split "`n" | Where-Object { $_ -match '^[A-Za-z0-9_-]+' } | Select-Object -First 1).Trim()
if (-not $compilerName) {
  Write-Error "Could not determine a compiler name from list output."
  exit 1
}

Write-Host "Setting default compiler to: $compilerName"
$setRes = Invoke-WSL @($CompactPath, 'config', 'set', 'default', $compilerName)
Write-Host $setRes

Write-Host "Running compile via npm..."
Push-Location (Join-Path $PSScriptRoot '..')
npm run compile
Pop-Location

