#!/bin/bash

# Create output directory
mkdir -p out

# Build the application
cd cmd/jarvis
go build -o ../../out/jarvis-mcp
cd ../..

# Make executable
chmod +x ./out/jarvis-mcp

echo "Build complete. Binary available at ./out/jarvis-mcp"
