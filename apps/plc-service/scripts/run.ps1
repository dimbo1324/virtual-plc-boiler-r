$PSScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location "$PSScriptRoot\.."

Write-Host "Current Directory: $(Get-Location)" -ForegroundColor Yellow
Write-Host "--- RUNNING GO TESTS ---" -ForegroundColor Cyan

go test ./... -v -cover

if ($LASTEXITCODE -eq 0) {
    Write-Host "`nSUCCESS: All tests passed!" -ForegroundColor Green
} else {
    Write-Host "`nWARNING: Some tests failed." -ForegroundColor Red
}

Write-Host "`nPress Enter to close..."
Read-Host