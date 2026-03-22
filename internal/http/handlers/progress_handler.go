package handlers

import (
	"interslavic/internal/models"
	"interslavic/internal/usecases"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProgressHandler struct {
	progressUseCase *usecases.ProgressUseCase
	validate        *validator.Validate
}

func NewProgressHandler(progressUseCase *usecases.ProgressUseCase) *ProgressHandler {
	return &ProgressHandler{
		progressUseCase: progressUseCase,
		validate:        validator.New(),
	}
}

// UpdateProgress godoc
// @Summary Update lesson progress
// @Description Update user's progress for a specific lesson (not_started, in_progress, completed)
// @Tags progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProgressRequest true "Progress update data"
// @Success 200 {object} models.LessonProgress
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/progress/update [post]
func (h *ProgressHandler) UpdateProgress(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста (устанавливается JWT middleware)
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	var req models.UpdateProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Валидация
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	progress, err := h.progressUseCase.UpdateProgress(c.UserContext(), userID, &req)
	if err != nil {
		if err.Error() == "lesson not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(progress)
}

// GetUserProgress godoc
// @Summary Get user's all progress
// @Description Get all progress records for the authenticated user
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.LessonProgress
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/progress [get]
func (h *ProgressHandler) GetUserProgress(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	progress, err := h.progressUseCase.GetUserProgress(c.UserContext(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(progress)
}

// GetLessonProgress godoc
// @Summary Get lesson progress
// @Description Get user's progress for a specific lesson
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Param lesson_id path int true "Lesson ID"
// @Success 200 {object} models.LessonProgress
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/progress/lesson/{lesson_id} [get]
func (h *ProgressHandler) GetLessonProgress(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	lessonID, err := strconv.Atoi(c.Params("lesson_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid lesson id",
		})
	}

	progress, err := h.progressUseCase.GetLessonProgress(c.UserContext(), userID, lessonID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(progress)
}

// GetCourseProgress godoc
// @Summary Get course progress
// @Description Get user's progress for a specific course with statistics
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "Course ID"
// @Success 200 {object} models.CourseProgress
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/progress/course/{course_id} [get]
func (h *ProgressHandler) GetCourseProgress(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	courseID, err := strconv.Atoi(c.Params("course_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid course id",
		})
	}

	progress, err := h.progressUseCase.GetCourseProgress(c.UserContext(), userID, courseID)
	if err != nil {
		if err.Error() == "course not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(progress)
}

// GetAllCoursesProgress godoc
// @Summary Get all courses progress
// @Description Get user's progress for all courses
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.CourseProgress
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/progress/courses [get]
func (h *ProgressHandler) GetAllCoursesProgress(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	progress, err := h.progressUseCase.GetAllCoursesProgress(c.UserContext(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(progress)
}

// GetCourseWithProgress godoc
// @Summary Get course with lessons progress
// @Description Get full course details with all lessons and user's progress for each lesson
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "Course ID"
// @Success 200 {object} models.CourseWithProgress
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/progress/course/{course_id}/full [get]
func (h *ProgressHandler) GetCourseWithProgress(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	courseID, err := strconv.Atoi(c.Params("course_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid course id",
		})
	}

	courseWithProgress, err := h.progressUseCase.GetCourseWithProgress(c.UserContext(), userID, courseID)
	if err != nil {
		if err.Error() == "course not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(courseWithProgress)
}
