#!/bin/bash

echo "🚀 Deploying GProc Backend to Fly.io..."

# Install flyctl if not present
if ! command -v flyctl &> /dev/null; then
    echo "Installing flyctl..."
    curl -L https://fly.io/install.sh | sh
    export PATH="$HOME/.fly/bin:$PATH"
fi

# Login to Fly.io (if not already logged in)
flyctl auth whoami || flyctl auth login

# Deploy backend
echo "Deploying backend..."
flyctl deploy --config fly.toml

echo "✅ Backend deployed to https://gproc-backend.fly.dev"