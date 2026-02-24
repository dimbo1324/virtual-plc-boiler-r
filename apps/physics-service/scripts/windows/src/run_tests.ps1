$Host.UI.RawUI.WindowTitle = "Physics Service - TESTS"
$Host.UI.RawUI.BackgroundColor = "DarkBlue"
$Host.UI.RawUI.ForegroundColor = "White"
Clear-Host
Write-Host "--- RUNNING TESTS ---" -ForegroundColor Cyan
$TimeStamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$OutDir = "out\tests"
$ReportFile = "$OutDir\report_$TimeStamp.html"
if (-not (Test-Path $OutDir)) {
    New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
}
try {
    Write-Host "Report path: $ReportFile" -ForegroundColor Yellow
    uv run pytest --html="$ReportFile" --self-contained-html
    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: All tests passed!" -ForegroundColor Green
    } else {
        Write-Host "WARNING: Some tests failed." -ForegroundColor Red
    }
}
catch {
    Write-Error "Error: $_"
}
Write-Host "Press Enter to close..."
Read-Host