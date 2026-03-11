#!/bin/bash
set -e

# Change directory to the root of the project
cd "$(dirname "$0")/.."

if command -v docker-compose >/dev/null 2>&1; then
	COMPOSE_CMD=(docker-compose)
elif command -v docker >/dev/null 2>&1; then
	COMPOSE_CMD=(docker compose)
else
	echo "Error: docker (or docker-compose) is not installed or not on PATH."
	exit 1
fi

echo "Stopping Docker containers..."
"${COMPOSE_CMD[@]}" down

echo ""
echo "Containers stopped successfully."
