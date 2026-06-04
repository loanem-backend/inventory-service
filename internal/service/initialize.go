package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loanem-backend/inventory-service/infra/database/sqlc"
	"github.com/loanem-backend/inventory-service/internal/repository"
)

func Initialize(p *pgxpool.Pool) (InstrumentService, ToolkitService, CourseService) {
	queries := sqlc.New(p)

	var (
		instrumentRepo = repository.NewInstrumentRepository(queries)
		toolkitRepo    = repository.NewToolkitRepository(queries)
		courseRepo     = repository.NewCourseRepository(queries)
	)

	var (
		instrumentServ = NewInstrumentService(instrumentRepo)
		toolkitServ    = NewToolkitService(toolkitRepo)
		courseServ     = NewCourseService(courseRepo)
	)

	return instrumentServ, toolkitServ, courseServ
}
