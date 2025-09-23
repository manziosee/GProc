#!/bin/bash

echo "🚀 Quick GProc Deployment"

# Build backend locally first
echo "Building backend..."
go build -o gproc cmd/main.go cmd/daemon.go

# Start backend locally for testing
echo "Starting backend on port 8080..."
./gproc daemon &
BACKEND_PID=$!

echo "Backend started with PID: $BACKEND_PID"
echo "Backend available at: http://localhost:8080"

# Build frontend
echo "Building frontend..."
cd fn
npm install
npm run build

# Serve frontend locally
echo "Starting frontend on port 5173..."
npm run preview &
FRONTEND_PID=$!

echo "Frontend started with PID: $FRONTEND_PID"
echo "Frontend available at: http://localhost:4173"

echo ""
echo "✅ GProc is now running locally:"
echo "   Frontend: http://localhost:4173"
echo "   Backend:  http://localhost:8080"
echo "   API:      http://localhost:8080/api/v1"
echo ""
echo "To stop services:"
echo "   kill $BACKEND_PID $FRONTEND_PID"