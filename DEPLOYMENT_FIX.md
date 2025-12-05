# Deployment Fix - Project Management System

## Issue

The deployment failed with the error:
```
npm error enoent Could not read package.json: Error: ENOENT: no such file or directory, open '/home/project/package.json'
```

## Root Cause

Memos is a hybrid Go/TypeScript application with a specific build structure:

1. Frontend code is in the `web/` subdirectory
2. The project uses **pnpm** (not npm) for frontend dependencies
3. There was no `package.json` at the project root for build systems to detect

## Solution Applied

I've added the following files to fix the deployment:

### 1. `/package.json` (Project Root)

A root-level package.json that orchestrates the build process:
- Installs pnpm
- Builds frontend from `web/` directory
- Builds Go backend
- Provides standard npm script interface for build systems

### 2. `/build.sh`

An executable build script that:
- Installs pnpm globally
- Builds frontend with proper output directory
- Compiles Go backend
- Works with any deployment platform

### 3. `/render.yaml`

Platform-specific configuration for Render.com that:
- Specifies Go environment
- Uses the build script
- Configures environment variables
- Sets up persistent storage

### 4. `/.npmrc`

Configuration to handle peer dependencies properly.

## How to Deploy

### Option 1: Automated (Recommended)

The build system should now automatically detect `package.json` and run:

```bash
npm run build
```

This will:
1. Install pnpm
2. Install frontend dependencies
3. Build frontend assets
4. Compile Go backend

### Option 2: Manual Build

Run the build script directly:

```bash
./build.sh
```

### Option 3: Docker

```bash
# Build frontend
cd web && pnpm install && pnpm run release && cd ..

# Build Docker image
docker build -f scripts/Dockerfile -t memos .

# Run container
docker run -p 5230:5230 -v ~/.memos/:/var/opt/memos memos
```

## Important: Protocol Buffers

⚠️ **CRITICAL**: The protobuf definitions have been updated for the new Area and Folder services, but the generated Go files have NOT been regenerated yet.

Before the backend can compile, you MUST run:

```bash
cd proto
buf generate
cd ..
```

This will generate:
- `proto/gen/api/v1/area_service.pb.go`
- `proto/gen/api/v1/area_service_grpc.pb.go`
- `proto/gen/api/v1/area_service.pb.gw.go`
- `proto/gen/api/v1/folder_service.pb.go`
- `proto/gen/api/v1/folder_service_grpc.pb.go`
- `proto/gen/api/v1/folder_service.pb.gw.go`
- Updated `proto/gen/api/v1/memo_service.pb.go` (with folder/area fields)

Without these files, the Go backend will fail to compile with import errors.

## Environment Variables

For production deployment:

```bash
MEMOS_MODE=prod
MEMOS_PORT=5230
MEMOS_DATA=/path/to/data/directory
```

## Database Migrations

The new project management system requires database migrations. The application will automatically run migrations on startup from:

- `store/migration/sqlite/0.26/00__project_management.sql`
- `store/migration/mysql/0.26/00__project_management.sql`
- `store/migration/postgres/0.26/00__project_management.sql`

## Verification

After deployment, verify the system is working:

1. **Check backend health**:
   ```bash
   curl http://your-domain/api/v1/health
   ```

2. **Test area creation**:
   ```bash
   curl -X POST http://your-domain/api/v1/areas \
     -H "Content-Type: application/json" \
     -d '{"area":{"displayName":"Test Area","description":"Test"}}'
   ```

3. **Open the web interface** and check for the Project Explorer in the sidebar

## Next Steps

After successful deployment:

1. **Regenerate protobuf files** (if not done automatically)
2. **Test the new features**:
   - Create areas
   - Create folders within areas
   - Assign memos to folders
3. **Integrate the UI components** into your navigation
4. **Review the documentation** in `PROJECT_MANAGEMENT_SYSTEM.md`

## Troubleshooting

### Build still fails with npm error

Make sure:
- The `package.json` is at the project root
- The build command is `npm run build` or `./build.sh`
- The deployment platform has access to install pnpm

### Go build fails with import errors

You need to regenerate protobuf files:
```bash
cd proto && buf generate && cd ..
```

### Frontend build fails

Delete and reinstall dependencies:
```bash
cd web
rm -rf node_modules .pnpm-store
pnpm install --frozen-lockfile
pnpm run release
```

## Support

If issues persist:
1. Check `BUILD_INSTRUCTIONS.md` for detailed build steps
2. Review `PROJECT_MANAGEMENT_SYSTEM.md` for feature documentation
3. Ensure all prerequisites are installed (Go 1.21+, Node 18+, pnpm 8+)
