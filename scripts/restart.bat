@echo off
REM Change to the root directory of the project
cd /d "%~dp0.."

echo ======================================
echo Restarting the Application...
echo ======================================

REM Stop running containers
call scripts\stop.bat

REM Start containers (rebuilding images if necessary)
call scripts\run.bat
