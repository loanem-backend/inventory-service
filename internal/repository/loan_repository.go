package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/entity"
	"github.com/oklog/ulid/v2"
)

type LoanRepository interface {
	Insert(ctx context.Context, l *entity.Loan) error

	FindByDate(ctx context.Context, date time.Time) ([]*entity.Loan, error)
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

func (r *loanRepository) FindByDate(ctx context.Context, date time.Time) ([]*entity.Loan, error) {
	rows, err := r.db.FindLoansByDate(ctx, pgtype.Date{Time: date, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Loan{}, nil
		}

		return nil, err
	}

	loans := make([]*entity.Loan, len(rows))

	for i, row := range rows {
		loans[i] = toLoan(sqlc.FindLoanByIDRow(row))
	}

	return loans, nil
}

func toLoan(row sqlc.FindLoanByIDRow) *entity.Loan {
	return &entity.Loan{
		ID: ulid.MustParse(row.ID),
		Toolkit: entity.Toolkit{
			ID:      int(row.ToolkitID.Int16),
			KitName: row.KitName.String,
			Course: entity.Course{
				ID:   int(row.CourseID.Int32),
				Name: row.CourseName.String,
			},
		},
		Team: entity.Team{
			ID: row.TeamID.String,
		},
		Submitter: entity.Participant{
			ID: row.SubmitterID,
		},
		Date:          row.Date.Time,
		SessionNumber: int(row.SessionNumber.Int16),
		Status:        entity.LoanStatus(row.Status.LoanStatus),
		Note:          row.Note.String,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
