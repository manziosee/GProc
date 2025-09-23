#!/bin/bash

echo "🚀 Deploying GProc Frontend to Fly.io..."

# Install flyctl if not present
if ! command -v flyctl &> /dev/null; then
    echo "Installing flyctl..."
    curl -L https://fly.io/install.sh | sh
    export PATH="$HOME/.fly/bin:$PATH"
fi

# Login to Fly.io (if not already logged in)
flyctl auth whoami || flyctl auth login

# Deploy frontend
echo "Deploying frontend..."
flyctl deploy --config fly.toml

echo "✅ Frontend deployed to https://gproc-frontend.fly.dev"