$projectRoot = (Get-Item $PSScriptRoot).Parent.Parent.Parent
Set-Location $projectRoot
uv run cmd/main.py
Write-Host "Project execution completed."