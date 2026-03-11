@echo off
REM Change to the root directory of the project
cd /d "%~dp0.."

echo ======================================
echo Stopping the Application...
echo ======================================

REM Stop and remove containers
docker compose down

echo.
echo ======================================
echo Application stopped successfully!
echo ======================================
