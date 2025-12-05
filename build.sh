#!/bin/bash
set -e

echo "==================================="
echo "Building Memos with Project Management"
echo "==================================="

# Ensure we're in the project root
cd "$(dirname "$0")"

echo ""
echo "Step 1: Installing pnpm globally..."
npm install -g pnpm@8.15.0

echo ""
echo "Step 2: Installing frontend dependencies..."
cd web
pnpm install --frozen-lockfile

echo ""
echo "Step 3: Building frontend..."
pnpm run release

echo ""
echo "Step 4: Building backend..."
cd ..
mkdir -p build
go build -o build/memos ./cmd/memos

echo ""
echo "==================================="
echo "Build completed successfully!"
echo "==================================="
echo ""
echo "Generated protobuf files need to be regenerated:"
echo "  cd proto && buf generate"
echo ""
echo "To run the application:"
echo "  ./build/memos --mode dev"
echo ""
