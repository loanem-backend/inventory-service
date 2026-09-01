package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/entity"
)

type LoanRepository interface {
	Insert(ctx context.Context, l *entity.Loan) error
}

type loanRepository struct {
	db *sqlc.Queries
}

func NewLoanRepository(q *sqlc.Queries) LoanRepository {
	return &loanRepository{
		db: q,
	}
}

func (s *loanRepository) Insert(ctx context.Context, l *entity.Loan) error {
	if err := s.db.InsertLoan(ctx, sqlc.InsertLoanParams{
		ID:            l.ID.String(),
		ToolkitID:     pgtype.Int2{Int16: int16(l.Toolkit.ID), Valid: true},
		TeamID:        pgtype.Text{String: l.Team.ID, Valid: l.Team.ID != ""},
		SubmitterID:   l.Submitter.ID,
		Date:          pgtype.Date{Time: l.Date, Valid: true},
		SessionNumber: pgtype.Int2{Int16: int16(l.SessionNumber), Valid: true},
		Status:        sqlc.NullLoanStatus{LoanStatus: sqlc.LoanStatus(l.Status), Valid: true},
	}); err != nil {
		return err
	}

	return nil
}
