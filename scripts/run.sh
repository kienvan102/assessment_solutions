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

echo "Building and starting Docker containers..."
"${COMPOSE_CMD[@]}" up -d --build

echo ""
echo "Project is running!"
echo "Web interface: http://localhost:3000"
echo "Backend API: http://localhost:8080"
echo "To view logs, run: ${COMPOSE_CMD[*]} logs -f"
