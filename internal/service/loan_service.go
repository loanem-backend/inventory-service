package service

import (
	"context"

	"github.com/loanem-backend/inventory-service/internal/entity"
	"github.com/loanem-backend/inventory-service/internal/repository"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoanService interface {
	CreateLoan(ctx context.Context, l *entity.Loan) (string, error)
}

type loanService struct {
	loanRepo repository.LoanRepository
}

func NewLoanService(lr repository.LoanRepository) LoanService {
	return &loanService{
		loanRepo: lr,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, l *entity.Loan) (string, error) {
	l.ID = ulid.Make()
	l.Status = entity.LoanStatusUpcoming

	if err := s.loanRepo.Insert(ctx, l); err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return l.ID.String(), nil
}
