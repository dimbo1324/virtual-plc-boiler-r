$projectRoot = (Get-Item $PSScriptRoot).Parent.Parent.Parent
Set-Location $projectRoot
uv run pytest
Write-Host "Tests execution completed."