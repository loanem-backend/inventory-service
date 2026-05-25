package repository

import (
	"context"

	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
)

type InstrumentRepository interface {
	Insert(ctx context.Context, name string) (int16, error)
}

type instrumentRepository struct {
	db *sqlc.Queries
}

func NewInstrumentRepository(q *sqlc.Queries) InstrumentRepository {
	return &instrumentRepository{
		db: q,
	}
}

func (r *instrumentRepository) Insert(ctx context.Context, name string) (int16, error) {
	result, err := r.db.InsertInstrument(ctx, name)
	if err != nil {
		return 0, err
	}

	return result, nil
}
