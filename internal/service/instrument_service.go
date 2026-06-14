package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/loanem-backend/api-gateway/pkg/storage"
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
	storage        *storage.S3Client
}

func NewInstrumentService(ir repository.InstrumentRepository, sc *storage.S3Client) InstrumentService {
	return &instrumentService{
		instrumentRepo: ir,
		storage:        sc,
	}
}

func (s *instrumentService) AddInstrument(ctx context.Context, name string) (int32, error) {
	instrumentID, err := s.instrumentRepo.Insert(ctx, &entity.Instrument{
		Name:    name,
		Picture: defaultInstrumentPicture,
	})
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

	if err := s.setInstrumentsPicture(ctx, instruments...); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return instruments, nil
}

func (s *instrumentService) setInstrumentsPicture(ctx context.Context, instruments ...*entity.Instrument) error {
	for _, i := range instruments {
		req, err := s.storage.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.storage.Bucket),
			Key:    aws.String(i.Picture),
		}, s3.WithPresignExpires(1*time.Hour))
		if err != nil {
			return fmt.Errorf("presign get object: %w", err)
		}

		i.Picture = req.URL
	}

	return nil
}

func (s *instrumentService) SetInstrumentPicture(ctx context.Context, instrument *entity.Instrument) error {
	instrument.UpdatedAt = time.Now()

	if instrument.Picture == "" {
		instrument.Picture = defaultInstrumentPicture
	}

	if err := s.instrumentRepo.UpdatePicture(ctx, instrument); err != nil {
		return err
	}

	return nil
}
