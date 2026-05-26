package repository

import (
	"context"

	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/entity"
)

type ToolkitRepository interface {
	Insert(ctx context.Context, t *entity.Toolkit) (int16, error)
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
