#!/bin/bash
# Development startup script for Entorno35
# Starts backend and frontend servers with proper environment configuration

BACKEND_PORT=8080
FRONTEND_PORT=3000

# Store the project root directory
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

echo "Starting Entorno35 Development Environment..."
echo "=============================================="
echo "Project root: $PROJECT_ROOT"

# Function to kill process on a specific port
kill_port() {
    local port=$1
    local pid=$(lsof -ti :$port 2>/dev/null)
    if [ -n "$pid" ]; then
        echo "Port $port is in use by PID $pid. Terminating..."
        kill -9 $pid 2>/dev/null || true
        sleep 1
    fi
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker Desktop."
    exit 1
fi

# Clean up ports that may be in use
echo ""
echo "Step 0: Checking for processes on required ports..."
kill_port $BACKEND_PORT
kill_port $FRONTEND_PORT
echo "Ports $BACKEND_PORT and $FRONTEND_PORT are now available."

# Start infrastructure
echo ""
echo "Step 1: Starting PostgreSQL and Redis..."
docker-compose -f "$PROJECT_ROOT/docker-compose.yml" up -d

# Wait for PostgreSQL to be ready
echo ""
echo "Step 2: Waiting for PostgreSQL to be ready..."
sleep 5

# Set environment variables for backend
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=entorno35
export DB_PASSWORD=entorno35
export DB_NAME=entorno35
export DB_SSLMODE=disable
export JWT_SECRET=your-secret-key-min-32-chars-long-for-development
export JWT_EXPIRY=24h
export CORS_ORIGIN=http://localhost:3000
export PORT=8080
export ENV=development

echo ""
echo "Step 3: Starting Backend API Server..."
echo "Backend will run on http://localhost:$BACKEND_PORT"
echo ""

# Start backend in background using subshell with absolute path
(cd "$PROJECT_ROOT/cmd/api" && go run main.go) &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

echo ""
echo "Step 4: Starting Frontend Development Server..."
echo "Frontend will run on http://localhost:$FRONTEND_PORT"
echo ""

# Start frontend in background using subshell with absolute path
(cd "$PROJECT_ROOT/web/frontend" && npm run dev) &
FRONTEND_PID=$!

echo ""
echo "=============================================="
echo "Entorno35 Development Environment Started!"
echo "=============================================="
echo ""
echo "Backend API: http://localhost:$BACKEND_PORT"
echo "Frontend UI: http://localhost:$FRONTEND_PORT"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "Stopping services..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
    sleep 1
    # Also kill any remaining processes on the ports
    kill_port $BACKEND_PORT
    kill_port $FRONTEND_PORT
    docker-compose -f "$PROJECT_ROOT/docker-compose.yml" down
    echo "All services stopped."
    exit 0
}

# Wait for user interrupt
trap cleanup INT TERM

# Keep script running
wait
