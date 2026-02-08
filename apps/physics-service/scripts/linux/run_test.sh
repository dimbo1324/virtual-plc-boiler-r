#!/bin/bash
set -e

echo "Running Tests..."
uv run pytest tests/ -v