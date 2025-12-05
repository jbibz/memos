#!/bin/bash
# Script to generate protobuf files
# This should be run before building the backend

set -e

echo "Checking for protobuf compiler..."

# Check if buf is installed
if ! command -v buf &> /dev/null; then
    echo "WARNING: 'buf' command not found."
    echo "Protobuf files need to be generated manually or buf needs to be installed."
    echo "The backend build will likely fail without generated protobuf files."
    echo ""
    echo "To install buf, visit: https://buf.build/docs/installation"
    echo ""
    echo "For now, skipping protobuf generation..."
    exit 0
fi

echo "Generating protobuf files..."
cd proto
buf generate
cd ..

echo "✓ Protobuf files generated successfully!"
