#!/usr/bin/env bash
set -e

# Change to the root directory of the project
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$ROOT_DIR"

echo "======================================"
echo "Running Backend Unit Tests..."
echo "======================================"

cd backend

# Ensure go modules are downloaded
go mod download

# Run tests recursively through all backend packages
# -v: verbose output to see which tests are running
# -count=1: bypass the test cache to ensure they always run
go test -v -count=1 ./...

echo "======================================"
echo "Backend Tests Completed Successfully!"
echo "======================================"
