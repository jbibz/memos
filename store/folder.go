package store

import (
	"context"
	"errors"

	"github.com/usememos/memos/internal/base"
)

type Folder struct {
	ID          int32
	UID         string
	RowStatus   RowStatus
	CreatorID   int32
	AreaID      int32
	CreatedTs   int64
	UpdatedTs   int64
	Name        string
	Description string
	ParentID    *int32
}

type FindFolder struct {
	ID        *int32
	UID       *string
	IDList    []int32
	UIDList   []string
	RowStatus *RowStatus
	CreatorID *int32
	AreaID    *int32
	ParentID  *int32
	Limit     *int
	Offset    *int
}

type UpdateFolder struct {
	ID          int32
	UID         *string
	UpdatedTs   *int64
	RowStatus   *RowStatus
	AreaID      *int32
	Name        *string
	Description *string
	ParentID    *int32
}

type DeleteFolder struct {
	ID int32
}

func (s *Store) CreateFolder(ctx context.Context, create *Folder) (*Folder, error) {
	if !base.UIDMatcher.MatchString(create.UID) {
		return nil, errors.New("invalid uid")
	}
	return s.driver.CreateFolder(ctx, create)
}

func (s *Store) ListFolders(ctx context.Context, find *FindFolder) ([]*Folder, error) {
	return s.driver.ListFolders(ctx, find)
}

func (s *Store) GetFolder(ctx context.Context, find *FindFolder) (*Folder, error) {
	list, err := s.ListFolders(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (s *Store) UpdateFolder(ctx context.Context, update *UpdateFolder) error {
	if update.UID != nil && !base.UIDMatcher.MatchString(*update.UID) {
		return errors.New("invalid uid")
	}
	return s.driver.UpdateFolder(ctx, update)
}

func (s *Store) DeleteFolder(ctx context.Context, delete *DeleteFolder) error {
	return s.driver.DeleteFolder(ctx, delete)
}
