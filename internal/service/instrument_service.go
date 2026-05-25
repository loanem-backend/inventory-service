package service

import (
	"context"

	"github.com/loanem-backend/inventory-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InstrumentService interface {
	AddInstrument(ctx context.Context, name string) (int32, error)
}

type instrumentService struct {
	instrumentRepo repository.InstrumentRepository
}

func NewInstrumentService(ir repository.InstrumentRepository) InstrumentService {
	return &instrumentService{
		instrumentRepo: ir,
	}
}

func (s *instrumentService) AddInstrument(ctx context.Context, name string) (int32, error) {
	instrumentID, err := s.instrumentRepo.Insert(ctx, name)
	if err != nil {
		return 0, status.Error(codes.Internal, "failed inserting instrument to database")
	}

	return int32(instrumentID), nil
}
