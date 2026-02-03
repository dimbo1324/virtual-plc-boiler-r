$srcPath = Join-Path $PSScriptRoot "src"
$scripts = Get-ChildItem -Path $srcPath -Filter *.ps1
foreach ($script in $scripts) {
    Start-Process powershell -ArgumentList "-NoExit", "-File $($script.FullName)"
}
Write-Host "All scripts from src launched in separate terminals."