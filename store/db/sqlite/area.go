package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/usememos/memos/store"
)

func (d *DB) CreateArea(ctx context.Context, create *store.Area) (*store.Area, error) {
	fields := []string{"`uid`", "`creator_id`", "`name`", "`description`"}
	placeholders := []string{"?", "?", "?", "?"}
	args := []any{create.UID, create.CreatorID, create.Name, create.Description}

	if create.ParentID != nil {
		fields = append(fields, "`parent_id`")
		placeholders = append(placeholders, "?")
		args = append(args, *create.ParentID)
	}

	query := `INSERT INTO area (` + strings.Join(fields, ", ") + `)
		VALUES (` + strings.Join(placeholders, ", ") + `)
		RETURNING id, uid, creator_id, created_ts, updated_ts, row_status, name, description, parent_id`

	var parentID sql.NullInt32
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(
		&create.ID,
		&create.UID,
		&create.CreatorID,
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

	area := create
	return area, nil
}

func (d *DB) ListAreas(ctx context.Context, find *store.FindArea) ([]*store.Area, error) {
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
	if v := find.ParentID; v != nil {
		if *v == 0 {
			where = append(where, "parent_id IS NULL")
		} else {
			where, args = append(where, "parent_id = ?"), append(args, *v)
		}
	}

	query := `SELECT id, uid, creator_id, created_ts, updated_ts, row_status, name, description, parent_id
		FROM area
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

	var areas []*store.Area
	for rows.Next() {
		area := &store.Area{}
		var parentID sql.NullInt32

		if err := rows.Scan(
			&area.ID,
			&area.UID,
			&area.CreatorID,
			&area.CreatedTs,
			&area.UpdatedTs,
			&area.RowStatus,
			&area.Name,
			&area.Description,
			&parentID,
		); err != nil {
			return nil, err
		}

		if parentID.Valid {
			area.ParentID = &parentID.Int32
		}

		areas = append(areas, area)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return areas, nil
}

func (d *DB) UpdateArea(ctx context.Context, update *store.UpdateArea) error {
	set, args := []string{}, []any{}

	if v := update.UpdatedTs; v != nil {
		set, args = append(set, "updated_ts = ?"), append(args, *v)
	} else {
		set, args = append(set, "updated_ts = ?"), append(args, time.Now().Unix())
	}
	if v := update.RowStatus; v != nil {
		set, args = append(set, "row_status = ?"), append(args, *v)
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
	query := `UPDATE area SET ` + strings.Join(set, ", ") + ` WHERE id = ?`
	if _, err := d.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteArea(ctx context.Context, delete *store.DeleteArea) error {
	query := `DELETE FROM area WHERE id = ?`
	if _, err := d.db.ExecContext(ctx, query, delete.ID); err != nil {
		return err
	}
	return nil
}
