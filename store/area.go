package store

import (
	"context"
	"errors"

	"github.com/usememos/memos/internal/base"
)

type Area struct {
	ID          int32
	UID         string
	RowStatus   RowStatus
	CreatorID   int32
	CreatedTs   int64
	UpdatedTs   int64
	Name        string
	Description string
	ParentID    *int32
}

type FindArea struct {
	ID        *int32
	UID       *string
	IDList    []int32
	UIDList   []string
	RowStatus *RowStatus
	CreatorID *int32
	ParentID  *int32
	Limit     *int
	Offset    *int
}

type UpdateArea struct {
	ID          int32
	UID         *string
	UpdatedTs   *int64
	RowStatus   *RowStatus
	Name        *string
	Description *string
	ParentID    *int32
}

type DeleteArea struct {
	ID int32
}

func (s *Store) CreateArea(ctx context.Context, create *Area) (*Area, error) {
	if !base.UIDMatcher.MatchString(create.UID) {
		return nil, errors.New("invalid uid")
	}
	return s.driver.CreateArea(ctx, create)
}

func (s *Store) ListAreas(ctx context.Context, find *FindArea) ([]*Area, error) {
	return s.driver.ListAreas(ctx, find)
}

func (s *Store) GetArea(ctx context.Context, find *FindArea) (*Area, error) {
	list, err := s.ListAreas(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (s *Store) UpdateArea(ctx context.Context, update *UpdateArea) error {
	if update.UID != nil && !base.UIDMatcher.MatchString(*update.UID) {
		return errors.New("invalid uid")
	}
	return s.driver.UpdateArea(ctx, update)
}

func (s *Store) DeleteArea(ctx context.Context, delete *DeleteArea) error {
	return s.driver.DeleteArea(ctx, delete)
}
