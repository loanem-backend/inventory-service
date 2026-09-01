package service

import (
	"context"
	"time"

	"github.com/loanem-backend/inventory-service/internal/entity"
	"github.com/loanem-backend/inventory-service/internal/repository"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoanService interface {
	Create(ctx context.Context, l *entity.Loan) (string, error)

	GetByID(ctx context.Context, loanID ulid.ULID) (*entity.Loan, error)
	GetByDate(ctx context.Context, date time.Time) ([]*entity.Loan, error)
}

type loanService struct {
	loanRepo repository.LoanRepository
}

func NewLoanService(lr repository.LoanRepository) LoanService {
	return &loanService{
		loanRepo: lr,
	}
}

func (s *loanService) Create(ctx context.Context, l *entity.Loan) (string, error) {
	l.ID = ulid.Make()
	l.Status = entity.LoanStatusUpcoming

	if err := s.loanRepo.Insert(ctx, l); err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return l.ID.String(), nil
}

func (s *loanService) GetByDate(ctx context.Context, date time.Time) ([]*entity.Loan, error) {
	loans, err := s.loanRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return loans, nil
}

func (s *loanService) GetByID(ctx context.Context, loanID ulid.ULID) (*entity.Loan, error) {
	loan, err := s.loanRepo.FindByID(ctx, loanID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return loan, nil
}
