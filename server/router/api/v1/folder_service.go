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

func (s *APIV1Service) CreateFolder(ctx context.Context, request *v1pb.CreateFolderRequest) (*v1pb.Folder, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user: %v", err)
	}

	folder := request.Folder
	if folder.DisplayName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "folder display name is required")
	}
	if folder.Area == "" {
		return nil, status.Errorf(codes.InvalidArgument, "folder area is required")
	}

	areaID, err := ExtractAreaIDFromName(folder.Area)
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

	uid, err := base.GenerateRandomID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate folder uid: %v", err)
	}

	create := &store.Folder{
		UID:         uid,
		CreatorID:   user.ID,
		AreaID:      area.ID,
		Name:        folder.DisplayName,
		Description: folder.Description,
	}

	if folder.Parent != "" {
		parentFolderID, err := ExtractFolderIDFromName(folder.Parent)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent folder name: %v", err)
		}
		parentFolder, err := s.Store.GetFolder(ctx, &store.FindFolder{
			UID: &parentFolderID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get parent folder: %v", err)
		}
		if parentFolder == nil {
			return nil, status.Errorf(codes.NotFound, "parent folder not found")
		}
		create.ParentID = &parentFolder.ID
	}

	createdFolder, err := s.Store.CreateFolder(ctx, create)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create folder: %v", err)
	}

	return convertFolderFromStore(createdFolder, area), nil
}

func (s *APIV1Service) ListFolders(ctx context.Context, request *v1pb.ListFoldersRequest) (*v1pb.ListFoldersResponse, error) {
	user, err := s.GetCurrentUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user: %v", err)
	}

	find := &store.FindFolder{
		CreatorID: &user.ID,
	}

	if request.PageSize > 0 {
		limit := int(request.PageSize)
		find.Limit = &limit
	}

	folders, err := s.Store.ListFolders(ctx, find)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list folders: %v", err)
	}

	response := &v1pb.ListFoldersResponse{}
	for _, folder := range folders {
		area, err := s.Store.GetArea(ctx, &store.FindArea{
			ID: &folder.AreaID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get area: %v", err)
		}
		response.Folders = append(response.Folders, convertFolderFromStore(folder, area))
	}

	return response, nil
}

func (s *APIV1Service) GetFolder(ctx context.Context, request *v1pb.GetFolderRequest) (*v1pb.Folder, error) {
	folderID, err := ExtractFolderIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid folder name: %v", err)
	}

	folder, err := s.Store.GetFolder(ctx, &store.FindFolder{
		UID: &folderID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get folder: %v", err)
	}
	if folder == nil {
		return nil, status.Errorf(codes.NotFound, "folder not found")
	}

	area, err := s.Store.GetArea(ctx, &store.FindArea{
		ID: &folder.AreaID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get area: %v", err)
	}

	return convertFolderFromStore(folder, area), nil
}

func (s *APIV1Service) UpdateFolder(ctx context.Context, request *v1pb.UpdateFolderRequest) (*v1pb.Folder, error) {
	folderID, err := ExtractFolderIDFromName(request.Folder.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid folder name: %v", err)
	}

	currentFolder, err := s.Store.GetFolder(ctx, &store.FindFolder{
		UID: &folderID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get folder: %v", err)
	}
	if currentFolder == nil {
		return nil, status.Errorf(codes.NotFound, "folder not found")
	}

	update := &store.UpdateFolder{
		ID: currentFolder.ID,
	}
	currentTime := time.Now().Unix()
	update.UpdatedTs = &currentTime

	if request.Folder.DisplayName != "" {
		update.Name = &request.Folder.DisplayName
	}
	if request.Folder.Description != "" {
		update.Description = &request.Folder.Description
	}

	if err := s.Store.UpdateFolder(ctx, update); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update folder: %v", err)
	}

	updatedFolder, err := s.Store.GetFolder(ctx, &store.FindFolder{
		ID: &currentFolder.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated folder: %v", err)
	}

	area, err := s.Store.GetArea(ctx, &store.FindArea{
		ID: &updatedFolder.AreaID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get area: %v", err)
	}

	return convertFolderFromStore(updatedFolder, area), nil
}

func (s *APIV1Service) DeleteFolder(ctx context.Context, request *v1pb.DeleteFolderRequest) (*emptypb.Empty, error) {
	folderID, err := ExtractFolderIDFromName(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid folder name: %v", err)
	}

	folder, err := s.Store.GetFolder(ctx, &store.FindFolder{
		UID: &folderID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get folder: %v", err)
	}
	if folder == nil {
		return nil, status.Errorf(codes.NotFound, "folder not found")
	}

	if err := s.Store.DeleteFolder(ctx, &store.DeleteFolder{
		ID: folder.ID,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete folder: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func convertFolderFromStore(folder *store.Folder, area *store.Area) *v1pb.Folder {
	result := &v1pb.Folder{
		Name:        fmt.Sprintf("folders/%s", folder.UID),
		Uid:         folder.UID,
		DisplayName: folder.Name,
		Description: folder.Description,
		Area:        fmt.Sprintf("areas/%s", area.UID),
		CreateTime:  timestamppb.New(time.Unix(folder.CreatedTs, 0)),
		UpdateTime:  timestamppb.New(time.Unix(folder.UpdatedTs, 0)),
	}

	if folder.ParentID != nil {
		result.Parent = fmt.Sprintf("folders/%d", *folder.ParentID)
	}

	return result
}

func ExtractFolderIDFromName(name string) (string, error) {
	return base.ExtractUIDFromResourceName(name, "folders")
}
