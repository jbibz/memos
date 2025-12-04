/*
  # Add Project Management System

  1. New Tables
    - `area`
      - `id` (integer, primary key, auto-increment)
      - `uid` (text, unique identifier)
      - `creator_id` (integer, foreign key to user)
      - `name` (text, area name)
      - `description` (text, optional description)
      - `created_ts` (bigint, creation timestamp)
      - `updated_ts` (bigint, last update timestamp)
      - `row_status` (text, NORMAL or ARCHIVED)
      - `parent_id` (integer, nullable, for nested areas)

    - `folder`
      - `id` (integer, primary key, auto-increment)
      - `uid` (text, unique identifier)
      - `creator_id` (integer, foreign key to user)
      - `area_id` (integer, foreign key to area)
      - `name` (text, folder name)
      - `description` (text, optional description)
      - `created_ts` (bigint, creation timestamp)
      - `updated_ts` (bigint, last update timestamp)
      - `row_status` (text, NORMAL or ARCHIVED)
      - `parent_id` (integer, nullable, for nested folders)

  2. Changes
    - Add `folder_id` column to `memo` table to link memos to folders
    - Add `area_id` column to `memo` table to link memos directly to areas (when not in a folder)

  3. Indexes
    - Index on `area.creator_id` for fast lookup
    - Index on `folder.creator_id` and `folder.area_id` for fast lookup
    - Index on `memo.folder_id` and `memo.area_id` for fast lookup

  4. Important Notes
    - Areas can be used for high-level organization (e.g., "Homelabs", "Work", "Personal")
    - Folders provide sub-organization within areas (e.g., "Equipment", "Self-hosted apps")
    - Memos can be linked to folders (folder_id) or directly to areas (area_id)
    - Both parent_id fields support future nested hierarchies
*/

-- Create area table
CREATE TABLE IF NOT EXISTS area (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uid TEXT NOT NULL UNIQUE,
  creator_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',
  name TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  parent_id INTEGER DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_area_creator_id ON area (creator_id);
CREATE INDEX IF NOT EXISTS idx_area_parent_id ON area (parent_id);

-- Create folder table
CREATE TABLE IF NOT EXISTS folder (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uid TEXT NOT NULL UNIQUE,
  creator_id INTEGER NOT NULL,
  area_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',
  name TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  parent_id INTEGER DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_folder_creator_id ON folder (creator_id);
CREATE INDEX IF NOT EXISTS idx_folder_area_id ON folder (area_id);
CREATE INDEX IF NOT EXISTS idx_folder_parent_id ON folder (parent_id);

-- Add folder_id and area_id columns to memo table
-- SQLite allows ALTER TABLE ADD COLUMN with IF NOT EXISTS from version 3.35.0+
-- However, for maximum compatibility, we'll use a simpler approach
ALTER TABLE memo ADD COLUMN folder_id INTEGER DEFAULT NULL;
ALTER TABLE memo ADD COLUMN area_id INTEGER DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_memo_folder_id ON memo (folder_id);
CREATE INDEX IF NOT EXISTS idx_memo_area_id ON memo (area_id);
