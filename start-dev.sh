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

# Function to kill all processes on a specific port (more thorough)
kill_port() {
    local port=$1
    # Try multiple methods to find and kill processes on the port
    # Method 1: lsof (standard)
    local pids=$(lsof -ti tcp:$port 2>/dev/null)
    if [ -n "$pids" ]; then
        echo "Port $port is in use by PID(s): $pids. Terminating..."
        echo "$pids" | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
    # Method 2: Check again with different lsof syntax
    pids=$(lsof -i :$port -t 2>/dev/null)
    if [ -n "$pids" ]; then
        echo "Port $port still in use by PID(s): $pids. Force terminating..."
        echo "$pids" | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# Function to wait for port to be available
wait_for_port_free() {
    local port=$1
    local max_attempts=5
    local attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if ! lsof -i :$port > /dev/null 2>&1; then
            return 0
        fi
        echo "Waiting for port $port to be free..."
        kill_port $port
        sleep 1
        attempt=$((attempt + 1))
    done
    echo "Warning: Port $port may still be in use"
    return 1
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
wait_for_port_free $BACKEND_PORT
wait_for_port_free $FRONTEND_PORT
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
export ENV=development
# Note: PORT is set per-process to avoid conflicts between backend (8080) and frontend (3000)

echo ""
echo "Step 3: Starting Backend API Server..."
echo "Backend will run on http://localhost:$BACKEND_PORT"
echo ""

# Start backend in background using subshell with absolute path
# PORT is set explicitly for the backend only to avoid conflicts with frontend
(cd "$PROJECT_ROOT/cmd/api" && PORT=$BACKEND_PORT exec go run main.go) &
BACKEND_PID=$!

# Wait for backend to start and verify it's running
echo "Waiting for backend to initialize..."
sleep 5

# Check if backend is responding
if curl -s --max-time 2 http://localhost:$BACKEND_PORT/health > /dev/null 2>&1 || \
   curl -s --max-time 2 http://localhost:$BACKEND_PORT > /dev/null 2>&1; then
    echo "Backend is running!"
else
    echo "Backend may still be starting..."
fi

echo ""
echo "Step 4: Starting Frontend Development Server..."
echo "Frontend will run on http://localhost:$FRONTEND_PORT"
echo ""

# Start frontend in background using subshell with absolute path
# Explicitly set port to avoid any conflicts
(cd "$PROJECT_ROOT/web/frontend" && exec npm run dev -- -p $FRONTEND_PORT) &
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
