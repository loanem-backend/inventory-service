package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/entity"
)

type ToolkitRepository interface {
	Insert(ctx context.Context, t *entity.Toolkit) (int16, error)
	FindAll(ctx context.Context) ([]*entity.Toolkit, error)
	FindWithInstruments(ctx context.Context) ([]*entity.Toolkit, error)
}

type toolkitRepository struct {
	db *sqlc.Queries
}

func NewToolkitRepository(q *sqlc.Queries) ToolkitRepository {
	return &toolkitRepository{
		db: q,
	}
}

func (r *toolkitRepository) Insert(ctx context.Context, t *entity.Toolkit) (int16, error) {
	result, err := r.db.InsertToolkit(ctx, sqlc.InsertToolkitParams{
		KitName:    t.KitName,
		TotalCount: int32(t.TotalCount),
	})
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *toolkitRepository) FindAll(ctx context.Context) ([]*entity.Toolkit, error) {
	rows, err := r.db.FindAllToolkits(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Toolkit{}, nil
		}

		return nil, err
	}

	toolkits := make([]*entity.Toolkit, len(rows))

	for i, row := range rows {
		toolkits[i] = toToolkit(row)
	}

	return toolkits, nil
}

func toToolkit(row sqlc.Toolkit) *entity.Toolkit {
	return &entity.Toolkit{
		ID:      int(row.ID),
		KitName: row.KitName,
		Course: entity.Course{
			ID: int(row.CourseID.Int32),
		},
		TotalCount:      int(row.TotalCount),
		OutOfOrderCount: int(row.OutOfOrderCount),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

func (r *toolkitRepository) FindWithInstruments(ctx context.Context) ([]*entity.Toolkit, error) {
	rows, err := r.db.FindToolkitsWithInstruments(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Toolkit{}, nil
		}

		return nil, err
	}

	toolkits := make([]*entity.Toolkit, 0)
	toolkitSeen := make(map[int16]int)
	toolkitIdx := 0

	for _, row := range rows {
		if _, seen := toolkitSeen[row.ToolkitID]; !seen {
			toolkits = append(toolkits, toToolkitFromWithInstruments(row))
			toolkits[toolkitIdx].Instruments = make([]entity.Instrument, 0)

			toolkitSeen[row.ToolkitID] = toolkitIdx
			toolkitIdx++
		}

		toolkits[toolkitSeen[row.ToolkitID]].Instruments = append(
			toolkits[toolkitSeen[row.ToolkitID]].Instruments,
			toInstrumentFromToolkitsWithInstruments(row),
		)
	}

	return toolkits, nil
}

func toToolkitFromWithInstruments(row sqlc.FindToolkitsWithInstrumentsRow) *entity.Toolkit {
	return &entity.Toolkit{
		ID:      int(row.ToolkitID),
		KitName: row.KitName,
		Course: entity.Course{
			ID: int(row.CourseID.Int32),
		},
		TotalCount:      int(row.TotalCount),
		OutOfOrderCount: int(row.OutOfOrderCount),
		CreatedAt:       row.ToolkitCreatedAt.Time,
		UpdatedAt:       row.ToolkitUpdatedAt.Time,
	}
}

func toInstrumentFromToolkitsWithInstruments(row sqlc.FindToolkitsWithInstrumentsRow) entity.Instrument {
	return entity.Instrument{
		ID:        int(row.InstrumentID),
		Name:      row.InstrumentName,
		CreatedAt: row.InstrumentCreatedAt.Time,
		UpdatedAt: row.InstrumentUpdatedAt.Time,
	}
}
