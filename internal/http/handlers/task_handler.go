package handlers

import (
	"interslavic/internal/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskUseCase *usecases.TaskUseCase
}

func NewTaskHandler(taskUseCase *usecases.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

// GetLessonTasks godoc
// @Summary Get tasks by lesson ID
// @Description Get all tasks of a specific lesson
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {array} models.Task
// @Failure 500 {object} map[string]string
// @Router /api/lessons/{id}/tasks [get]
func (h *TaskHandler) GetLessonTasks(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid lesson id",
		})
	}

	tasks, err := h.taskUseCase.GetTasksByLessonID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(tasks)
}

// CheckTaskAnswer godoc
// @Summary Check task answer
// @Description Check if the provided answer is correct
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.CheckAnswerRequest true "Answer to check"
// @Success 200 {object} usecases.CheckAnswerResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tasks/check [post]
func (h *TaskHandler) CheckTaskAnswer(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	var req usecases.CheckAnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.TaskID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "task_id is required",
		})
	}

	// Используем userID для логирования, но пока не сохраняем результат
	_ = userID

	resp, err := h.taskUseCase.CheckAnswer(c.UserContext(), req.TaskID, req.Answer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}
