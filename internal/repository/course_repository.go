package repository

import (
	"context"

	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/entity"
)

type CourseRepository interface {
	Insert(ctx context.Context, c *entity.Course) error
}

type courseRepository struct {
	db *sqlc.Queries
}

func NewCourseRepository(q *sqlc.Queries) CourseRepository {
	return &courseRepository{
		db: q,
	}
}

func (r *courseRepository) Insert(ctx context.Context, c *entity.Course) error {
	if err := r.db.InsertCourse(ctx, sqlc.InsertCourseParams{
		ID:   int32(c.ID),
		Name: c.Name,
		Year: int32(c.Year),
	}); err != nil {
		return err
	}

	return nil
}
