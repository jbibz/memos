/*
  # Add Project Management System

  1. New Tables
    - `area`
      - `id` (integer, primary key, auto-increment)
      - `uid` (varchar, unique identifier)
      - `creator_id` (integer, foreign key to user)
      - `name` (text, area name)
      - `description` (text, optional description)
      - `created_ts` (bigint, creation timestamp)
      - `updated_ts` (bigint, last update timestamp)
      - `row_status` (varchar, NORMAL or ARCHIVED)
      - `parent_id` (integer, nullable, for nested areas)

    - `folder`
      - `id` (integer, primary key, auto-increment)
      - `uid` (varchar, unique identifier)
      - `creator_id` (integer, foreign key to user)
      - `area_id` (integer, foreign key to area)
      - `name` (text, folder name)
      - `description` (text, optional description)
      - `created_ts` (bigint, creation timestamp)
      - `updated_ts` (bigint, last update timestamp)
      - `row_status` (varchar, NORMAL or ARCHIVED)
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
CREATE TABLE IF NOT EXISTS `area` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `uid` VARCHAR(255) NOT NULL UNIQUE,
  `creator_id` INT NOT NULL,
  `created_ts` BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `updated_ts` BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `row_status` VARCHAR(32) NOT NULL DEFAULT 'NORMAL' CHECK (`row_status` IN ('NORMAL', 'ARCHIVED')),
  `name` TEXT NOT NULL,
  `description` TEXT NOT NULL,
  `parent_id` INT DEFAULT NULL,
  KEY `idx_area_creator_id` (`creator_id`),
  KEY `idx_area_parent_id` (`parent_id`)
);

-- Create folder table
CREATE TABLE IF NOT EXISTS `folder` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `uid` VARCHAR(255) NOT NULL UNIQUE,
  `creator_id` INT NOT NULL,
  `area_id` INT NOT NULL,
  `created_ts` BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `updated_ts` BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
  `row_status` VARCHAR(32) NOT NULL DEFAULT 'NORMAL' CHECK (`row_status` IN ('NORMAL', 'ARCHIVED')),
  `name` TEXT NOT NULL,
  `description` TEXT NOT NULL,
  `parent_id` INT DEFAULT NULL,
  KEY `idx_folder_creator_id` (`creator_id`),
  KEY `idx_folder_area_id` (`area_id`),
  KEY `idx_folder_parent_id` (`parent_id`)
);

-- Add folder_id and area_id columns to memo table
ALTER TABLE `memo` ADD COLUMN IF NOT EXISTS `folder_id` INT DEFAULT NULL;
ALTER TABLE `memo` ADD COLUMN IF NOT EXISTS `area_id` INT DEFAULT NULL;

CREATE INDEX IF NOT EXISTS `idx_memo_folder_id` ON `memo` (`folder_id`);
CREATE INDEX IF NOT EXISTS `idx_memo_area_id` ON `memo` (`area_id`);
