package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/repository"
	"github.com/loanem-backend/inventory-service/internal/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/grpc"
)

func Start(s *grpc.Server, p *pgxpool.Pool) {
	queries := sqlc.New(p)

	var (
		instrumentRepo = repository.NewInstrumentRepository(queries)
	)

	var (
		instrumentServ = service.NewInstrumentService(instrumentRepo)
	)

	pbinventory.RegisterInventoryServiceServer(s, NewInventoryServer(
		NewInstrumentServer(instrumentServ),
	))
}
