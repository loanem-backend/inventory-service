package service

import (
	"context"

	"github.com/loanem-backend/inventory-service/internal/entity"
	"github.com/loanem-backend/inventory-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ToolkitService interface {
	Create(ctx context.Context, t *entity.Toolkit) (int16, error)
}

type toolkitService struct {
	toolkitRepo repository.ToolkitRepository
}

func NewToolkitService(tr repository.ToolkitRepository) ToolkitService {
	return &toolkitService{
		toolkitRepo: tr,
	}
}

func (s *toolkitService) Create(ctx context.Context, t *entity.Toolkit) (int16, error) {
	toolkitID, err := s.toolkitRepo.Insert(ctx, t)
	if err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}

	return toolkitID, nil
}
