package server

import (
	"github.com/loanem-backend/inventory-service/internal/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/grpc"
)

func Start(s *grpc.Server, is service.InstrumentService, ts service.ToolkitService) {
	pbinventory.RegisterInstrumentServiceServer(s, NewInstrumentServer(is))
	pbinventory.RegisterToolkitServiceServer(s, NewToolkitServer(ts))
}
