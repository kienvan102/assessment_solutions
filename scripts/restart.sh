#!/usr/bin/env bash
set -e

# Change to the root directory of the project
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "======================================"
echo "Restarting the Application..."
echo "======================================"

# Stop running containers
./scripts/stop.sh

# Start containers (rebuilding images if necessary)
./scripts/run.sh
