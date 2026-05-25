package server

import (
	"context"

	"github.com/loanem-backend/inventory-service/internal/mapper"
	"github.com/loanem-backend/inventory-service/internal/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

type InstrumentServer struct {
	serv service.InstrumentService
}

func NewInstrumentServer(is service.InstrumentService) *InstrumentServer {
	return &InstrumentServer{
		serv: is,
	}
}

func (s *InstrumentServer) AddInstrument(ctx context.Context, req *pbinventory.AddInstrumentRequest) (*pbinventory.AddInstrumentResponse, error) {
	idData, err := s.serv.AddInstrument(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return mapper.IntToAddInstrumentResponse(idData), nil
}
