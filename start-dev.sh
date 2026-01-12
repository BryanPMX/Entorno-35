#!/bin/bash
# Development startup script for Entorno35
# Starts backend and frontend servers with proper environment configuration

set -e

echo "Starting Entorno35 Development Environment..."
echo "=============================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker Desktop."
    exit 1
fi

# Start infrastructure
echo ""
echo "Step 1: Starting PostgreSQL and Redis..."
make docker-up

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
echo "Backend will run on http://localhost:8080"
echo ""

# Start backend in background
cd cmd/api && go run main.go &
BACKEND_PID=$!
cd ../..

# Wait for backend to start
sleep 3

echo ""
echo "Step 4: Starting Frontend Development Server..."
echo "Frontend will run on http://localhost:3000"
echo ""

# Start frontend in background
cd web/frontend && npm run dev &
FRONTEND_PID=$!
cd ../..

echo ""
echo "=============================================="
echo "Entorno35 Development Environment Started!"
echo "=============================================="
echo ""
echo "Backend API: http://localhost:8080"
echo "Frontend UI: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for user interrupt
trap "echo 'Stopping services...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; make docker-down; exit 0" INT TERM

# Keep script running
wait
