package server

import (
	"context"

	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

type InventoryServer struct {
	pbinventory.UnimplementedInventoryServiceServer
	instrumentSrv InstrumentServer
}

func NewInventoryServer(is *InstrumentServer) *InventoryServer {
	return &InventoryServer{
		instrumentSrv: *is,
	}
}

func (s *InventoryServer) AddInstrument(ctx context.Context, req *pbinventory.AddInstrumentRequest) (*pbinventory.AddInstrumentResponse, error) {
	return s.instrumentSrv.AddInstrument(ctx, req)
}
