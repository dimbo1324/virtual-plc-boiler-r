$Host.UI.RawUI.WindowTitle = "Physics Service - MAIN APP"
$Host.UI.RawUI.BackgroundColor = "Black"
$Host.UI.RawUI.ForegroundColor = "Green"
Clear-Host
Write-Host "------------------------------------------------" -ForegroundColor Green
Write-Host "           STARTING MAIN SERVICE                " -ForegroundColor Green
Write-Host "------------------------------------------------"
if (-not (Get-Command "uv" -ErrorAction SilentlyContinue)) {
    Write-Error "UV not found! Please ensure it is installed and added to PATH."
    exit
}
try {
    Write-Host "Executing: uv run cmd/main.py" -ForegroundColor Gray
    uv run cmd/main.py
}
catch {
    Write-Error "Error while starting application: $_"
}
if ($LASTEXITCODE -ne 0) {
    Write-Host "`nProcess finished with exit code $LASTEXITCODE" -ForegroundColor Yellow
    Read-Host "Press Enter to close..."
}