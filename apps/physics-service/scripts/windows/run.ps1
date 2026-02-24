$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "Requesting Administrator privileges..." -ForegroundColor Yellow
    Start-Process powershell.exe -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`Detailed report will be in out/tests`"" -Verb RunAs
    Exit
}
$ScriptDir = $PSScriptRoot
$SrcDir = Join-Path $ScriptDir "src"
$ProjectRoot = Resolve-Path (Join-Path $ScriptDir "..\..")
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   PHYSICS SERVICE: STARTING STACK        " -ForegroundColor Cyan
Write-Host "=========================================="
Write-Host "Root: $ProjectRoot" -ForegroundColor Gray
if (Test-Path $SrcDir) {
    $Scripts = Get-ChildItem -Path $SrcDir -Filter "*.ps1"
    foreach ($Script in $Scripts) {
        Write-Host "Launching: $($Script.Name)..." -ForegroundColor Green
        Start-Process powershell.exe -ArgumentList "-NoExit", "-ExecutionPolicy Bypass", "-File `"$($Script.FullName)`"" -WorkingDirectory $ProjectRoot
    }
} else {
    Write-Error "Directory 'src' not found: $SrcDir"
}
Write-Host "Done. Windows are opening." -ForegroundColor Cyan
Start-Sleep -Seconds 2