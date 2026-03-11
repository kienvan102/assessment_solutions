# Assessment Solutions

This repository contains the solutions for the assessment. The solutions are provided via a Go backend API and displayed through a React + TypeScript frontend.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Project Structure

- `backend/` - Go API serving the solutions (requires manual `go mod init` and module setup)
- `web/` - React + TypeScript frontend
- `scripts/` - Utility scripts to run and manage the project

## Running the Project

The easiest way to start the project is using **Make** (recommended) or the platform-specific scripts.

### Using Make (Recommended - Works on All Platforms)

If you have `make` installed, you can use these simple commands:

```bash
make run       # Start the application
make stop      # Stop the application
make restart   # Restart the application
make test      # Run backend tests
make logs      # View container logs
make status    # Show running containers
make clean     # Stop and remove everything including volumes
make help      # Show all available commands
```

The Makefile automatically detects your OS and runs the appropriate script.

---

### Using Platform-Specific Scripts

If you prefer to use the scripts directly:

### For macOS/Linux Users

### 1. Start the Application

```bash
./scripts/run.sh
```
This script will build the Docker images and start the containers in detached mode.

### 2. View the Solutions

Once the containers are running, you can access the frontend web interface at:
**http://localhost:3000**

The Go backend API is accessible at:
**http://localhost:8080**

### 3. Stop the Application

To stop and remove the running containers, use:

```bash
./scripts/stop.sh
```

### 4. Restart the Application

To quickly stop, rebuild, and start the containers, use:

```bash
./scripts/restart.sh
```

### 5. Run Backend Tests

To run the Go backend unit tests:

```bash
./scripts/test.sh
```

---

### For Windows Users

### 1. Start the Application

```cmd
scripts\run.bat
```
This script will build the Docker images and start the containers in detached mode.

### 2. View the Solutions

Once the containers are running, you can access the frontend web interface at:
**http://localhost:3000**

The Go backend API is accessible at:
**http://localhost:8080**

### 3. Stop the Application

To stop and remove the running containers, use:

```cmd
scripts\stop.bat
```

### 4. Restart the Application

To quickly stop, rebuild, and start the containers, use:

```cmd
scripts\restart.bat
```

### 5. Run Backend Tests

To run the Go backend unit tests:

```cmd
scripts\test.bat
```

## CI/CD

This project includes automated testing via GitHub Actions. The CI pipeline runs automatically on:
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main`, `master`, or `develop` branches

### What the CI Pipeline Does:
- ✅ Runs all backend Go unit tests
- ✅ Checks for race conditions with `-race` flag
- ✅ Generates code coverage reports
- ✅ Uploads coverage artifacts (retained for 30 days)

### Viewing Test Results:
Once you push to GitHub, you can view test results in the **Actions** tab of your repository.

### Running Tests Locally:
You can run the same tests locally using:
```bash
make test
# or
./scripts/test.sh  # macOS/Linux
scripts\test.bat   # Windows
```

---

## Development

If you wish to run the components locally without Docker:

**Backend:**
Navigate to `backend/` and initialize the Go module:
```bash
cd backend
go mod init <module-name>
go mod tidy
go run main.go
```

**Frontend:**
Navigate to `web/` and start the Vite development server:
```bash
cd web
npm install
npm run dev
```
