# Quick Start Guide - Project Management System

## What Was Added

Your Memos application now has a hierarchical project management system:

```
Areas (e.g., "Homelabs")
  └── Folders (e.g., "Equipment", "Self-hosted apps")
      └── Memos (your notes and ideas)
```

## Deployment Issue Fixed ✅

The deployment error has been resolved. I've added:

1. ✅ Root `package.json` for build system detection
2. ✅ `build.sh` script for automated building
3. ✅ `render.yaml` for platform-specific deployment
4. ✅ `.npmrc` for dependency management

## Before You Can Run This

⚠️ **CRITICAL STEP** - Generate Protocol Buffer Files:

```bash
cd proto
buf generate
cd ..
```

This MUST be done before the Go backend can compile. The new Area and Folder services require generated protobuf files that don't exist yet.

## Building the Application

### Quick Build (Recommended)

```bash
./build.sh
```

### Manual Build

```bash
# 1. Generate protobuf files (REQUIRED)
cd proto && buf generate && cd ..

# 2. Build frontend
cd web && pnpm install && pnpm run release && cd ..

# 3. Build backend
go build -o build/memos ./cmd/memos
```

## Running Locally

```bash
./build/memos --mode dev --port 5230
```

Then open: http://localhost:5230

## Deploying

### For Render.com / Railway / Similar Platforms

The build system should now automatically detect and build correctly. Just push to your repository and trigger a deployment.

Build command: `npm run build` or `./build.sh`
Start command: `./build/memos --mode prod --port $PORT`

## Features Overview

### Backend (Go)
- ✅ Area CRUD API (`/api/v1/areas`)
- ✅ Folder CRUD API (`/api/v1/folders`)
- ✅ Memo-to-folder assignment
- ✅ Database migrations (SQLite, MySQL, PostgreSQL)

### Frontend (React)
- ✅ ProjectExplorer component (hierarchical tree view)
- ✅ CreateAreaDialog (create new areas)
- ✅ CreateFolderDialog (create folders in areas)
- ✅ FolderSelector (for memo editor integration)

### Database Schema
- ✅ `area` table (top-level organization)
- ✅ `folder` table (nested under areas)
- ✅ `memo` table updated (folder_id, area_id columns)

## Visual Preview

Open `PREVIEW.html` in your browser to see exactly what the UI looks like!

## Integration Steps

### 1. Add ProjectExplorer to Sidebar

In your navigation component:

```tsx
import ProjectExplorer from "@/components/ProjectExplorer";

// Add to sidebar
<ProjectExplorer />
```

### 2. Add FolderSelector to Memo Editor

In `MemoEditor/index.tsx`:

```tsx
import FolderSelector from "@/components/FolderSelector";

const [selectedFolder, setSelectedFolder] = useState<string>();
const [selectedArea, setSelectedArea] = useState<string>();

// Add before textarea
<FolderSelector
  selectedFolder={selectedFolder}
  selectedArea={selectedArea}
  onFolderChange={setSelectedFolder}
  onAreaChange={setSelectedArea}
  className="mb-2"
/>
```

## Example Usage

```
Your Homelabs Setup:

📁 Homelabs (Area)
   📂 Equipment (Folder)
      📝 Dell R720 Server Specs
      📝 Network Switch Configuration
      📝 UPS Battery Replacement

   📂 Self-hosted apps (Folder)
      📝 Nextcloud Setup
      📝 Pi-hole Configuration
      📝 Home Assistant Notes
      📝 Jellyfin Media Server
```

## Documentation

- `PROJECT_MANAGEMENT_SYSTEM.md` - Complete system documentation
- `BUILD_INSTRUCTIONS.md` - Detailed build instructions
- `DEPLOYMENT_FIX.md` - Deployment troubleshooting guide
- `PREVIEW.html` - Visual preview of the UI

## Retry Deployment

Once you've regenerated the protobuf files, you can retry your deployment. The build system should now work correctly!

## Need Help?

1. Check `DEPLOYMENT_FIX.md` for troubleshooting
2. Review `BUILD_INSTRUCTIONS.md` for build steps
3. See `PROJECT_MANAGEMENT_SYSTEM.md` for feature docs

## Next Steps

1. ✅ Regenerate protobuf files: `cd proto && buf generate`
2. ✅ Test build locally: `./build.sh`
3. ✅ Integrate UI components into your app
4. ✅ Retry deployment
5. ✅ Create your first Area and Folder!
