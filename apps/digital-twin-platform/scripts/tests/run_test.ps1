$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Definition
$projectRoot = Resolve-Path "$scriptPath\..\.."
$outDir = "$projectRoot\out"

if (-not (Test-Path $outDir)) {
    New-Item -ItemType Directory -Path $outDir | Out-Null
}

Set-Location $projectRoot

Write-Host "`n>>> Step 1: Running unit tests..." -ForegroundColor Cyan
go test -v ./...

if ($LASTEXITCODE -ne 0) {
    Write-Host "`n[!] Some tests failed. Coverage report will not be generated." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "`n>>> Step 2: Generating coverage profile in /out..." -ForegroundColor Cyan
$coveragePath = "$outDir\coverage.out"
go test -coverprofile="$coveragePath" ./...

if (Test-Path $coveragePath) {
    Write-Host "`n>>> Step 3: Opening HTML coverage report..." -ForegroundColor Green
    go tool cover -html="$coveragePath"
}

Write-Host "`nDone! All artifacts are in $outDir" -ForegroundColor Yellow