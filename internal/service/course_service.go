package service

import (
	"context"

	"github.com/loanem-backend/course-service/pkg/messaging"
	"github.com/loanem-backend/inventory-service/internal/entity"
	"github.com/loanem-backend/inventory-service/internal/repository"
)

type CourseService interface {
	Add(ctx context.Context, arg messaging.CourseEvent) error
	Remove(ctx context.Context, arg messaging.CourseEvent) error
}

type courseService struct {
	courseRepo repository.CourseRepository
}

func NewCourseService(cr repository.CourseRepository) CourseService {
	return &courseService{
		courseRepo: cr,
	}
}

func (s *courseService) Add(ctx context.Context, arg messaging.CourseEvent) error {
	if err := s.courseRepo.Insert(ctx, &entity.Course{
		ID:   arg.CourseID,
		Name: arg.CourseName,
		Year: arg.CourseYear,
	}); err != nil {
		return err
	}

	return nil
}

func (s *courseService) Remove(ctx context.Context, arg messaging.CourseEvent) error {
	if err := s.courseRepo.Delete(ctx, int32(arg.CourseID)); err != nil {
		return err
	}

	return nil
}
