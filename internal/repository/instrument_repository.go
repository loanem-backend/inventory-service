package repository

import (
	"context"

	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/entity"
)

type InstrumentRepository interface {
	Insert(ctx context.Context, name string) (int16, error)
	Delete(ctx context.Context, iID int16) error
	FindAll(ctx context.Context) ([]*entity.Instrument, error)
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

func (r *instrumentRepository) Delete(ctx context.Context, iID int16) error {
	if err := r.db.DeleteInstrumentByID(ctx, iID); err != nil {
		return err
	}

	return nil
}

func (r *instrumentRepository) FindAll(ctx context.Context) ([]*entity.Instrument, error) {
	rows, err := r.db.FindAllInstruments(ctx)
	if err != nil {
		return nil, err
	}

	instruments := make([]*entity.Instrument, len(rows))

	for i, row := range rows {
		instruments[i] = toInstrument(row)
	}

	return instruments, nil
}

func toInstrument(row sqlc.Instrument) *entity.Instrument {
	return &entity.Instrument{
		ID:        int(row.ID),
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
