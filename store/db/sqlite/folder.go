package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/usememos/memos/store"
)

func (d *DB) CreateFolder(ctx context.Context, create *store.Folder) (*store.Folder, error) {
	fields := []string{"`uid`", "`creator_id`", "`area_id`", "`name`", "`description`"}
	placeholders := []string{"?", "?", "?", "?", "?"}
	args := []any{create.UID, create.CreatorID, create.AreaID, create.Name, create.Description}

	if create.ParentID != nil {
		fields = append(fields, "`parent_id`")
		placeholders = append(placeholders, "?")
		args = append(args, *create.ParentID)
	}

	query := `INSERT INTO folder (` + strings.Join(fields, ", ") + `)
		VALUES (` + strings.Join(placeholders, ", ") + `)
		RETURNING id, uid, creator_id, area_id, created_ts, updated_ts, row_status, name, description, parent_id`

	var parentID sql.NullInt32
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&create.ID,
		&create.UID,
		&create.CreatorID,
		&create.AreaID,
		&create.CreatedTs,
		&create.UpdatedTs,
		&create.RowStatus,
		&create.Name,
		&create.Description,
		&parentID,
	); err != nil {
		return nil, err
	}

	if parentID.Valid {
		create.ParentID = &parentID.Int32
	}

	folder := create
	return folder, nil
}

func (d *DB) ListFolders(ctx context.Context, find *store.FindFolder) ([]*store.Folder, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.UID; v != nil {
		where, args = append(where, "uid = ?"), append(args, *v)
	}
	if v := find.IDList; len(v) > 0 {
		placeholders := make([]string, len(v))
		for i, id := range v {
			placeholders[i] = "?"
			args = append(args, id)
		}
		where = append(where, "id IN ("+strings.Join(placeholders, ",")+")")
	}
	if v := find.UIDList; len(v) > 0 {
		placeholders := make([]string, len(v))
		for i, uid := range v {
			placeholders[i] = "?"
			args = append(args, uid)
		}
		where = append(where, "uid IN ("+strings.Join(placeholders, ",")+")")
	}
	if v := find.RowStatus; v != nil {
		where, args = append(where, "row_status = ?"), append(args, *v)
	}
	if v := find.CreatorID; v != nil {
		where, args = append(where, "creator_id = ?"), append(args, *v)
	}
	if v := find.AreaID; v != nil {
		where, args = append(where, "area_id = ?"), append(args, *v)
	}
	if v := find.ParentID; v != nil {
		if *v == 0 {
			where = append(where, "parent_id IS NULL")
		} else {
			where, args = append(where, "parent_id = ?"), append(args, *v)
		}
	}

	query := `SELECT id, uid, creator_id, area_id, created_ts, updated_ts, row_status, name, description, parent_id
		FROM folder
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY created_ts DESC`

	if v := find.Limit; v != nil {
		query += " LIMIT ?"
		args = append(args, *v)
	}
	if v := find.Offset; v != nil {
		query += " OFFSET ?"
		args = append(args, *v)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*store.Folder
	for rows.Next() {
		folder := &store.Folder{}
		var parentID sql.NullInt32

		if err := rows.Scan(
			&folder.ID,
			&folder.UID,
			&folder.CreatorID,
			&folder.AreaID,
			&folder.CreatedTs,
			&folder.UpdatedTs,
			&folder.RowStatus,
			&folder.Name,
			&folder.Description,
			&parentID,
		); err != nil {
			return nil, err
		}

		if parentID.Valid {
			folder.ParentID = &parentID.Int32
		}

		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

func (d *DB) UpdateFolder(ctx context.Context, update *store.UpdateFolder) error {
	set, args := []string{}, []any{}

	if v := update.UpdatedTs; v != nil {
		set, args = append(set, "updated_ts = ?"), append(args, *v)
	} else {
		set, args = append(set, "updated_ts = ?"), append(args, time.Now().Unix())
	}
	if v := update.RowStatus; v != nil {
		set, args = append(set, "row_status = ?"), append(args, *v)
	}
	if v := update.AreaID; v != nil {
		set, args = append(set, "area_id = ?"), append(args, *v)
	}
	if v := update.Name; v != nil {
		set, args = append(set, "name = ?"), append(args, *v)
	}
	if v := update.Description; v != nil {
		set, args = append(set, "description = ?"), append(args, *v)
	}
	if v := update.ParentID; v != nil {
		set, args = append(set, "parent_id = ?"), append(args, *v)
	}
	if v := update.UID; v != nil {
		set, args = append(set, "uid = ?"), append(args, *v)
	}

	args = append(args, update.ID)
	query := `UPDATE folder SET ` + strings.Join(set, ", ") + ` WHERE id = ?`
	if _, err := d.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteFolder(ctx context.Context, delete *store.DeleteFolder) error {
	query := `DELETE FROM folder WHERE id = ?`
	if _, err := d.db.ExecContext(ctx, query, delete.ID); err != nil {
		return err
	}
	return nil
}
