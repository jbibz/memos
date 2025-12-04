package v1

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/usememos/memos/internal/base"
	"github.com/usememos/memos/store"

	v1pb "github.com/usememos/memos/proto/gen/api/v1"
)

func (s *APIV1Service) CreateArea(ctx context.Context, request *v1pb.CreateAreaRequest) (*v1pb.Area, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user: %v", err)
	}

	area := request.Area
	if area.DisplayName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "area display name is required")
	}

	uid, err := base.GenerateRandomID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate area uid: %v", err)
	}

	create := &store.Area{
		UID:         uid,
		CreatorID:   user.ID,
		Name:        area.DisplayName,
		Description: area.Description,
	}

	if area.Parent != "" {
		parentAreaID, err := ExtractAreaIDFromName(area.Parent)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent area name: %v", err)
		}
		parentID := int32(parentAreaID)
		create.ParentID = &parentID
	}

	createdArea, err := s.Store.CreateArea(ctx, create)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create area: %v", err)
	}

	return convertAreaFromStore(createdArea), nil
}

func (s *APIV1Service) ListAreas(ctx context.Context, request *v1pb.ListAreasRequest) (*v1pb.ListAreasResponse, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user: %v", err)
	}

	find := &store.FindArea{
		CreatorID: &user.ID,
	}

	if request.PageSize > 0 {
		limit := int(request.PageSize)
		find.Limit = &limit
	}

	areas, err := s.Store.ListAreas(ctx, find)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list areas: %v", err)
	}

	response := &v1pb.ListAreasResponse{}
	for _, area := range areas {
		response.Areas = append(response.Areas, convertAreaFromStore(area))
	}

	return response, nil
}

func (s *APIV1Service) GetArea(ctx context.Context, request *v1pb.GetAreaRequest) (*v1pb.Area, error) {
	areaID, err := ExtractAreaIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid area name: %v", err)
	}

	area, err := s.Store.GetArea(ctx, &store.FindArea{
		UID: &areaID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get area: %v", err)
	}
	if area == nil {
		return nil, status.Errorf(codes.NotFound, "area not found")
	}

	return convertAreaFromStore(area), nil
}

func (s *APIV1Service) UpdateArea(ctx context.Context, request *v1pb.UpdateAreaRequest) (*v1pb.Area, error) {
	areaID, err := ExtractAreaIDFromName(request.Area.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid area name: %v", err)
	}

	currentArea, err := s.Store.GetArea(ctx, &store.FindArea{
		UID: &areaID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get area: %v", err)
	}
	if currentArea == nil {
		return nil, status.Errorf(codes.NotFound, "area not found")
	}

	update := &store.UpdateArea{
		ID: currentArea.ID,
	}
	currentTime := time.Now().Unix()
	update.UpdatedTs = &currentTime

	if request.Area.DisplayName != "" {
		update.Name = &request.Area.DisplayName
	}
	if request.Area.Description != "" {
		update.Description = &request.Area.Description
	}

	if err := s.Store.UpdateArea(ctx, update); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update area: %v", err)
	}

	updatedArea, err := s.Store.GetArea(ctx, &store.FindArea{
		ID: &currentArea.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated area: %v", err)
	}

	return convertAreaFromStore(updatedArea), nil
}

func (s *APIV1Service) DeleteArea(ctx context.Context, request *v1pb.DeleteAreaRequest) (*emptypb.Empty, error) {
	areaID, err := ExtractAreaIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid area name: %v", err)
	}

	area, err := s.Store.GetArea(ctx, &store.FindArea{
		UID: &areaID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get area: %v", err)
	}
	if area == nil {
		return nil, status.Errorf(codes.NotFound, "area not found")
	}

	if err := s.Store.DeleteArea(ctx, &store.DeleteArea{
		ID: area.ID,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete area: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func convertAreaFromStore(area *store.Area) *v1pb.Area {
	result := &v1pb.Area{
		Name:        fmt.Sprintf("areas/%s", area.UID),
		Uid:         area.UID,
		DisplayName: area.Name,
		Description: area.Description,
		CreateTime:  timestamppb.New(time.Unix(area.CreatedTs, 0)),
		UpdateTime:  timestamppb.New(time.Unix(area.UpdatedTs, 0)),
	}

	if area.ParentID != nil {
		result.Parent = fmt.Sprintf("areas/%d", *area.ParentID)
	}

	return result
}

func ExtractAreaIDFromName(name string) (string, error) {
	return base.ExtractUIDFromResourceName(name, "areas")
}
