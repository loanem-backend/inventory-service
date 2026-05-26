package server

import (
	"context"

	"github.com/loanem-backend/inventory-service/internal/mapper"
	"github.com/loanem-backend/inventory-service/internal/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

type ToolkitServer struct {
	pbinventory.UnimplementedToolkitServiceServer
	toolkitServ service.ToolkitService
}

func NewToolkitServer(ts service.ToolkitService) *ToolkitServer {
	return &ToolkitServer{
		toolkitServ: ts,
	}
}

func (s *ToolkitServer) AddToolkit(ctx context.Context, req *pbinventory.AddToolkitRequest) (*pbinventory.AddToolkitResponse, error) {
	idData, err := s.toolkitServ.Create(ctx, mapper.AddToolkitRequestToToolkit(req))
	if err != nil {
		return nil, err
	}

	return mapper.IntToAddToolkitResponse(idData), nil
}
