# Build Instructions

This project is Memos with an added Project Management System (Areas and Folders).

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and pnpm 8+
- Protocol Buffers compiler (protoc) and buf

## Important Note About Protobuf

⚠️ **The protobuf definitions have been updated but the generated Go files have NOT been regenerated yet.**

Before building the backend, you MUST regenerate the protobuf files:

```bash
cd proto
buf generate
cd ..
```

## Build Process

### Option 1: Automated Build Script

The easiest way to build the entire project:

```bash
./build.sh
```

This will:
1. Install pnpm globally
2. Install frontend dependencies
3. Build the frontend
4. Build the backend

### Option 2: Manual Build

#### Step 1: Generate Protobuf Files (REQUIRED)

```bash
cd proto
buf generate
cd ..
```

#### Step 2: Build Frontend

```bash
cd web
pnpm install
pnpm run release
cd ..
```

#### Step 3: Build Backend

```bash
mkdir -p build
go build -o build/memos ./cmd/memos
```

## Running the Application

### Development Mode

```bash
./build/memos --mode dev --port 5230
```

### Production Mode

```bash
./build/memos --mode prod --port 5230
```

## Deployment

### Environment Variables

- `MEMOS_MODE`: Set to `prod` for production
- `MEMOS_PORT`: Port number (default: 5230)
- `MEMOS_DATA`: Data directory path

### Using Docker

Build the Docker image:

```bash
# Build frontend first
cd web && pnpm install && pnpm run release && cd ..

# Build Docker image
docker build -f scripts/Dockerfile -t memos-pm .
```

Run the container:

```bash
docker run -d \
  -p 5230:5230 \
  -v ~/.memos/:/var/opt/memos \
  --name memos-pm \
  memos-pm
```

## Project Management System

The new features include:
- Areas (top-level organization)
- Folders (sub-organization within areas)
- Hierarchical memo organization

See `PROJECT_MANAGEMENT_SYSTEM.md` for detailed documentation.

## Troubleshooting

### "npm not found" error

This project uses pnpm, not npm. Install pnpm:

```bash
npm install -g pnpm@8.15.0
```

### "package.json not found" error

Make sure you're running commands from the correct directory:
- Frontend commands: `web/` directory
- Backend commands: project root
- Build script: project root

### Missing protobuf files

If you see import errors related to `proto/gen/api/v1`, regenerate the protobuf files:

```bash
cd proto
buf generate
cd ..
```

### Frontend build fails

Delete node_modules and reinstall:

```bash
cd web
rm -rf node_modules pnpm-lock.yaml
pnpm install
pnpm run release
cd ..
```
