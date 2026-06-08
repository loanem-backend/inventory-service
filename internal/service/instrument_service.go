package service

import (
	"context"
	"time"

	"github.com/loanem-backend/inventory-service/internal/entity"
	"github.com/loanem-backend/inventory-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InstrumentService interface {
	AddInstrument(ctx context.Context, name string) (int32, error)
	RemoveInstrument(ctx context.Context, instrumentID int32) error
	GetAllInstruments(ctx context.Context) ([]*entity.Instrument, error)
	SetInstrumentPicture(ctx context.Context, instrument *entity.Instrument) error
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

func (s *instrumentService) RemoveInstrument(ctx context.Context, instrumentID int32) error {
	if err := s.instrumentRepo.Delete(ctx, int16(instrumentID)); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *instrumentService) GetAllInstruments(ctx context.Context) ([]*entity.Instrument, error) {
	instruments, err := s.instrumentRepo.FindAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// set picture url

	return instruments, nil
}

func (s *instrumentService) SetInstrumentPicture(ctx context.Context, instrument *entity.Instrument) error {
	instrument.UpdatedAt = time.Now()

	if err := s.instrumentRepo.UpdatePicture(ctx, instrument); err != nil {
		return err
	}

	return nil
}
