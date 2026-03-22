package usecases

import (
	"context"
	"errors"
	"interslavic/internal/database"
	"interslavic/internal/models"
	"interslavic/logging"
)

type CourseUseCase struct {
	courseRepo database.CourseRepository
	logger     *logging.ModuleLogger
}

func NewCourseUseCase(
	courseRepo database.CourseRepository,
	logger *logging.ModuleLogger,
) *CourseUseCase {
	return &CourseUseCase{
		courseRepo: courseRepo,
		logger:     logging.NewModuleLogger("COURSE", "USECASE", logger),
	}
}

func (uc *CourseUseCase) GetAllCourses(ctx context.Context) ([]models.Course, error) {
	courses, err := uc.courseRepo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("failed to get courses", logging.ErrAttr(err))
		return nil, errors.New("failed to get courses")
	}
	return courses, nil
}

func (uc *CourseUseCase) GetCourseByID(ctx context.Context, id int) (*models.Course, error) {
	course, err := uc.courseRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get course", logging.ErrAttr(err))
		return nil, errors.New("failed to get course")
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	return course, nil
}