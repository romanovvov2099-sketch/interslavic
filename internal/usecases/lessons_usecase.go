package usecases

import (
	"context"
	"errors"
	"interslavic/internal/database"
	"interslavic/internal/models"
	"interslavic/logging"
)

type LessonUseCase struct {
	lessonRepo database.LessonRepository
	taskRepo   database.TaskRepository
	logger     *logging.ModuleLogger
}

func NewLessonUseCase(
	lessonRepo database.LessonRepository,
	taskRepo database.TaskRepository,
	logger *logging.ModuleLogger,
) *LessonUseCase {
	return &LessonUseCase{
		lessonRepo: lessonRepo,
		taskRepo:   taskRepo,
		logger:     logging.NewModuleLogger("LESSON", "USECASE", logger),
	}
}

func (uc *LessonUseCase) GetLessonByID(ctx context.Context, id int) (*models.Lesson, error) {
	lesson, err := uc.lessonRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get lesson", logging.ErrAttr(err))
		return nil, errors.New("failed to get lesson")
	}
	if lesson == nil {
		return nil, errors.New("lesson not found")
	}
	return lesson, nil
}

func (uc *LessonUseCase) GetLessonsByCourseID(ctx context.Context, courseID int) ([]models.Lesson, error) {
	lessons, err := uc.lessonRepo.FindByCourseID(ctx, courseID)
	if err != nil {
		uc.logger.Error("failed to get lessons", logging.ErrAttr(err))
		return nil, errors.New("failed to get lessons")
	}
	return lessons, nil
}

type LessonWithTasks struct {
	Lesson models.Lesson `json:"lesson"`
	Tasks  []models.Task `json:"tasks"`
}

func (uc *LessonUseCase) GetLessonWithTasks(ctx context.Context, lessonID int) (*LessonWithTasks, error) {
	lesson, err := uc.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	tasks, err := uc.taskRepo.FindByLessonID(ctx, lessonID)
	if err != nil {
		uc.logger.Error("failed to get tasks", logging.ErrAttr(err))
		return nil, errors.New("failed to get tasks")
	}

	// Скрываем ответы для заданий с выбором ответа
	for i := range tasks {
		if tasks[i].TaskType == models.ChoisesTask {
			tasks[i].Answer = ""
		}
	}

	return &LessonWithTasks{
		Lesson: *lesson,
		Tasks:  tasks,
	}, nil
}