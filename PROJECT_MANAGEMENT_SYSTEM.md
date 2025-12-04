# Project Management System

This document describes the project management system that has been added to the Memos application, allowing hierarchical organization of memos into Areas and Folders.

## Overview

The project management system introduces a three-level hierarchy:

1. **Areas** - Top-level containers (e.g., "Homelabs", "Work", "Personal")
2. **Folders** - Sub-containers within areas (e.g., "Equipment", "Self-hosted apps")
3. **Memos** - Individual notes/ideas within folders or areas

## Architecture

### Database Schema

Three new tables have been created:

#### `area` table
- `id`: Auto-increment primary key
- `uid`: Unique identifier
- `creator_id`: Foreign key to user
- `name`: Area display name
- `description`: Optional description
- `created_ts`, `updated_ts`: Timestamps
- `row_status`: NORMAL or ARCHIVED
- `parent_id`: For future nested areas support

#### `folder` table
- `id`: Auto-increment primary key
- `uid`: Unique identifier
- `creator_id`: Foreign key to user
- `area_id`: Foreign key to area
- `name`: Folder display name
- `description`: Optional description
- `created_ts`, `updated_ts`: Timestamps
- `row_status`: NORMAL or ARCHIVED
- `parent_id`: For future nested folders support

#### `memo` table updates
Two new nullable columns added:
- `folder_id`: Links memo to a folder
- `area_id`: Links memo directly to an area (when not in a folder)

### Backend Implementation

#### Store Layer (`/store`)
- `area.go`: Area CRUD operations
- `folder.go`: Folder CRUD operations
- Database drivers implemented for SQLite (MySQL and PostgreSQL need implementation)

#### API Layer (`/server/router/api/v1`)
- `area_service.go`: gRPC service for area management
- `folder_service.go`: gRPC service for folder management
- Services registered in `v1.go`

#### Protocol Buffers (`/proto/api/v1`)
- `area_service.proto`: Area API definitions
- `folder_service.proto`: Folder API definitions
- `memo_service.proto`: Updated with folder and area fields

### Frontend Components

#### Core Components
1. **ProjectExplorer** (`/web/src/components/ProjectExplorer.tsx`)
   - Hierarchical tree view of areas and folders
   - Expand/collapse areas
   - Create new areas and folders
   - Click folders to filter memos

2. **CreateAreaDialog** (`/web/src/components/CreateAreaDialog.tsx`)
   - Dialog for creating new areas
   - Name and description fields

3. **CreateFolderDialog** (`/web/src/components/CreateFolderDialog.tsx`)
   - Dialog for creating new folders
   - Area selection dropdown
   - Name and description fields

4. **FolderSelector** (`/web/src/components/FolderSelector.tsx`)
   - Reusable component for selecting area and folder
   - Can be integrated into memo editor

#### Type Definitions
- `/web/src/types/proto/api/v1/area_service.ts`
- `/web/src/types/proto/api/v1/folder_service.ts`

## API Endpoints

### Areas
- `POST /api/v1/areas` - Create area
- `GET /api/v1/areas` - List areas
- `GET /api/v1/areas/{uid}` - Get area
- `PATCH /api/v1/areas/{uid}` - Update area
- `DELETE /api/v1/areas/{uid}` - Delete area

### Folders
- `POST /api/v1/folders` - Create folder
- `GET /api/v1/folders` - List folders
- `GET /api/v1/folders/{uid}` - Get folder
- `PATCH /api/v1/folders/{uid}` - Update folder
- `DELETE /api/v1/folders/{uid}` - Delete folder

## Integration Steps

### 1. Database Migration

The database migrations are located in:
- SQLite: `/store/migration/sqlite/0.26/00__project_management.sql`
- MySQL: `/store/migration/mysql/0.26/00__project_management.sql`
- PostgreSQL: `/store/migration/postgres/0.26/00__project_management.sql`

The LATEST.sql files have also been updated to include the new schema.

### 2. Backend Compilation

After the protobuf definitions are updated, you need to regenerate the Go code:

```bash
# From the project root
cd proto
buf generate
```

Then rebuild the application:

```bash
go build -o memos ./cmd/memos
```

### 3. Frontend Integration

#### Add ProjectExplorer to Navigation

In your main navigation component (e.g., `Navigation.tsx` or `NavigationDrawer.tsx`), add:

```tsx
import ProjectExplorer from "@/components/ProjectExplorer";

// Inside your navigation component
<ProjectExplorer />
```

#### Add FolderSelector to Memo Editor

In the MemoEditor component, add the FolderSelector:

```tsx
import FolderSelector from "@/components/FolderSelector";

// In the MemoEditor state, add:
const [selectedFolder, setSelectedFolder] = useState<string>();
const [selectedArea, setSelectedArea] = useState<string>();

// In the editor UI, add the selector:
<FolderSelector
  selectedFolder={selectedFolder}
  selectedArea={selectedArea}
  onFolderChange={setSelectedFolder}
  onAreaChange={setSelectedArea}
  className="mb-2"
/>

// When creating/updating a memo, include the folder/area:
const memo = {
  ...otherMemoFields,
  folder: selectedFolder,
  area: !selectedFolder ? selectedArea : undefined, // Only set area if no folder
};
```

### 4. Install Dependencies

The frontend uses these UI components from shadcn/ui:
- Button
- Dialog
- Input
- Label
- Select
- Textarea

These should already be available in the project.

## Usage Example

```typescript
// Example: Your Homelabs hierarchy
- Homelabs (Area)
  - Equipment (Folder)
    - "Dell R720 Server Specs" (Memo)
    - "Network Switch Configuration" (Memo)
  - Self-hosted apps (Folder)
    - "Nextcloud Setup Guide" (Memo)
    - "Pi-hole Configuration" (Memo)
    - "Home Assistant Notes" (Memo)
```

## Future Enhancements

1. **Nested Hierarchies**: Support for nested areas and folders using the `parent_id` field
2. **Drag and Drop**: Move memos between folders with drag-and-drop
3. **Bulk Operations**: Move multiple memos at once
4. **Folder Colors/Icons**: Customize folder appearance
5. **Filtering**: Filter memos by area/folder in the main view
6. **Search**: Search within specific areas/folders
7. **Permissions**: Share specific areas/folders with other users
8. **MySQL/PostgreSQL Drivers**: Implement the area/folder operations for these databases

## Notes

- The system is designed to be backward compatible - existing memos without folder/area assignments continue to work
- Areas and folders are user-specific - each user has their own organizational structure
- The `parent_id` fields in both tables are prepared for future nested hierarchy support
- Row-level security and proper authentication are handled by the existing user system
