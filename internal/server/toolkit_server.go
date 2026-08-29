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

func (s *ToolkitServer) GetAllToolkits(ctx context.Context, req *pbinventory.GetAllToolkitsRequest) (*pbinventory.GetAllToolkitsResponse, error) {
	toolkitsData, err := s.toolkitServ.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if toolkitsData == nil {
		return &pbinventory.GetAllToolkitsResponse{
			Toolkits: []*pbinventory.Toolkit{},
		}, nil
	}

	return mapper.ToolkitsToGetAllToolkitsResponse(toolkitsData), nil
}

func (s *ToolkitServer) GetToolkitsWithInstruments(ctx context.Context, req *pbinventory.GetToolkitsWithInstrumentsRequest) (*pbinventory.GetToolkitsWithInstrumentsResponse, error) {
	toolkitsData, err := s.toolkitServ.GetWithInstruments(ctx)
	if err != nil {
		return nil, err
	}

	if toolkitsData == nil {
		return &pbinventory.GetToolkitsWithInstrumentsResponse{
			Toolkits: []*pbinventory.ToolkitWithInstruments{},
		}, nil
	}

	return mapper.ToolkitsToGetToolkitsWithInstrumentsResponse(toolkitsData), nil
}
