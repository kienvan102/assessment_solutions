@echo off
REM Change to the root directory of the project
cd /d "%~dp0.."

echo ======================================
echo Running Backend Unit Tests...
echo ======================================

cd backend

REM Ensure go modules are downloaded
go mod download

REM Run tests recursively through all backend packages
REM -v: verbose output to see which tests are running
REM -count=1: bypass the test cache to ensure they always run
go test -v -count=1 ./...

echo.
echo ======================================
echo Backend Tests Completed Successfully!
echo ======================================
