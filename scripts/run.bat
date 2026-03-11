@echo off
REM Change to the root directory of the project
cd /d "%~dp0.."

echo ======================================
echo Starting the Application...
echo ======================================

REM Build and start containers in detached mode
docker compose up --build -d

echo.
echo ======================================
echo Application started successfully!
echo ======================================
echo Frontend: http://localhost:3000
echo Backend:  http://localhost:8080
echo.
