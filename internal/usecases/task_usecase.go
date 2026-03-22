package usecases

import (
	"context"
	"errors"
	"interslavic/internal/database"
	"interslavic/internal/models"
	"interslavic/logging"
)

type TaskUseCase struct {
	taskRepo database.TaskRepository
	logger   *logging.ModuleLogger
}

func NewTaskUseCase(
	taskRepo database.TaskRepository,
	logger *logging.ModuleLogger,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo: taskRepo,
		logger:   logging.NewModuleLogger("TASK", "USECASE", logger),
	}
}

func (uc *TaskUseCase) GetTasksByLessonID(ctx context.Context, lessonID int) ([]models.Task, error) {
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

	return tasks, nil
}

type CheckAnswerRequest struct {
	TaskID int    `json:"task_id"`
	Answer string `json:"answer"`
}

type CheckAnswerResponse struct {
	IsCorrect bool   `json:"is_correct"`
	Message   string `json:"message,omitempty"`
}

func (uc *TaskUseCase) CheckAnswer(ctx context.Context, taskID int, answer string) (*CheckAnswerResponse, error) {
	task, err := uc.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		uc.logger.Error("failed to find task", logging.ErrAttr(err))
		return nil, errors.New("task not found")
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	isCorrect, err := uc.taskRepo.CheckAnswer(ctx, taskID, answer)
	if err != nil {
		uc.logger.Error("failed to check answer", logging.ErrAttr(err))
		return nil, errors.New("failed to check answer")
	}

	response := &CheckAnswerResponse{
		IsCorrect: isCorrect,
	}

	if isCorrect {
		response.Message = "Correct answer!"
	} else {
		response.Message = "Incorrect answer. Try again!"
	}

	return response, nil
}