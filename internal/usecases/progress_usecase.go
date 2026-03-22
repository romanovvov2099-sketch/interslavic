package usecases

import (
	"context"
	"errors"
	"interslavic/internal/database"
	"interslavic/internal/models"
	"interslavic/logging"
	"time"
)

type ProgressUseCase struct {
	progressRepo database.LessonProgressRepository
	courseRepo   database.CourseRepository
	lessonRepo   database.LessonRepository
	logger       *logging.ModuleLogger
}

func NewProgressUseCase(
	progressRepo database.LessonProgressRepository,
	lessonRepo database.LessonRepository,
	logger *logging.ModuleLogger,
) *ProgressUseCase {
	return &ProgressUseCase{
		progressRepo: progressRepo,
		lessonRepo:   lessonRepo,
		logger:       logging.NewModuleLogger("PROGRESS", "USECASE", logger),
	}
}

func (uc *ProgressUseCase) UpdateProgress(ctx context.Context, userID int, req *models.UpdateProgressRequest) (*models.LessonProgress, error) {
	// Проверяем существование урока
	lesson, err := uc.lessonRepo.FindByID(ctx, req.LessonID)
	if err != nil {
		uc.logger.Error("failed to find lesson", logging.ErrAttr(err))
		return nil, errors.New("failed to get lesson")
	}
	if lesson == nil {
		return nil, errors.New("lesson not found")
	}

	// Создаем объект прогресса
	progress := &models.LessonProgress{
		UserID:   userID,
		LessonID: req.LessonID,
		Status:   req.Status,
		Score:    req.Score,
	}

	// Если статус "completed", устанавливаем дату завершения
	if req.Status == "completed" {
		now := time.Now()
		progress.CompletionDate = &now
	}

	// Сохраняем прогресс
	err = uc.progressRepo.CreateOrUpdate(ctx, progress)
	if err != nil {
		uc.logger.Error("failed to save progress", logging.ErrAttr(err))
		return nil, errors.New("failed to save progress")
	}

	return progress, nil
}

func (uc *ProgressUseCase) GetUserProgress(ctx context.Context, userID int) ([]models.LessonProgress, error) {
	progress, err := uc.progressRepo.FindByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get user progress", logging.ErrAttr(err))
		return nil, errors.New("failed to get progress")
	}
	return progress, nil
}

func (uc *ProgressUseCase) GetLessonProgress(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error) {
	progress, err := uc.progressRepo.FindByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// Возвращаем дефолтный прогресс, если записи нет
			return &models.LessonProgress{
				UserID:   userID,
				LessonID: lessonID,
				Status:   "not_started",
				Score:    0,
			}, nil
		}
		uc.logger.Error("failed to get lesson progress", logging.ErrAttr(err))
		return nil, errors.New("failed to get progress")
	}
	return progress, nil
}

// GetCourseProgress - получить прогресс по конкретному курсу
func (uc *ProgressUseCase) GetCourseProgress(ctx context.Context, userID, courseID int) (*models.CourseProgress, error) {
	// Проверяем существование курса
	course, err := uc.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		uc.logger.Error("failed to find course", logging.ErrAttr(err))
		return nil, errors.New("failed to get course")
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	progress, err := uc.progressRepo.GetCourseProgressStats(ctx, userID, courseID)
	if err != nil {
		uc.logger.Error("failed to get course progress", logging.ErrAttr(err))
		return nil, errors.New("failed to get course progress")
	}

	return progress, nil
}

// GetAllCoursesProgress - получить прогресс по всем курсам
func (uc *ProgressUseCase) GetAllCoursesProgress(ctx context.Context, userID int) ([]models.CourseProgress, error) {
	progress, err := uc.progressRepo.GetAllCoursesProgress(ctx, userID)
	if err != nil {
		uc.logger.Error("failed to get all courses progress", logging.ErrAttr(err))
		return nil, errors.New("failed to get courses progress")
	}

	return progress, nil
}

// GetCourseWithProgress - получить курс со списком уроков и прогрессом по каждому
func (uc *ProgressUseCase) GetCourseWithProgress(ctx context.Context, userID, courseID int) (*models.CourseWithProgress, error) {
	// Получаем информацию о курсе
	course, err := uc.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		uc.logger.Error("failed to find course", logging.ErrAttr(err))
		return nil, errors.New("failed to get course")
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	// Получаем уроки с прогрессом
	lessons, err := uc.progressRepo.GetLessonsWithProgressByCourse(ctx, userID, courseID)
	if err != nil {
		uc.logger.Error("failed to get lessons with progress", logging.ErrAttr(err))
		return nil, errors.New("failed to get lessons")
	}

	// Получаем статистику курса
	stats, err := uc.progressRepo.GetCourseProgressStats(ctx, userID, courseID)
	if err != nil {
		uc.logger.Error("failed to get course stats", logging.ErrAttr(err))
		return nil, errors.New("failed to get course stats")
	}

	result := &models.CourseWithProgress{
		Course:           *course,
		Lessons:          lessons,
		ProgressPercent:  stats.ProgressPercent,
		CompletedLessons: stats.CompletedLessons,
		TotalLessons:     stats.TotalLessons,
		Status:           stats.Status,
	}

	return result, nil
}
