package handlers

import (
	"interslavic/internal/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LessonHandler struct {
	lessonUseCase *usecases.LessonUseCase
}

func NewLessonHandler(lessonUseCase *usecases.LessonUseCase) *LessonHandler {
	return &LessonHandler{
		lessonUseCase: lessonUseCase,
	}
}

// GetLessonByID godoc
// @Summary Get lesson by ID
// @Description Get lesson details by ID
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} models.Lesson
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid lesson id",
		})
	}

	lesson, err := h.lessonUseCase.GetLessonByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(lesson)
}

// GetCourseLessons godoc
// @Summary Get lessons by course ID
// @Description Get all lessons of a specific course
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {array} models.Lesson
// @Failure 500 {object} map[string]string
// @Router /api/courses/{id}/lessons [get]
func (h *LessonHandler) GetCourseLessons(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid course id",
		})
	}

	lessons, err := h.lessonUseCase.GetLessonsByCourseID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(lessons)
}

// GetLessonWithTasks godoc
// @Summary Get lesson with tasks
// @Description Get lesson details and all its tasks
// @Tags lessons
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} usecases.LessonWithTasks
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/lessons/{id}/full [get]
func (h *LessonHandler) GetLessonWithTasks(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid lesson id",
		})
	}

	lessonWithTasks, err := h.lessonUseCase.GetLessonWithTasks(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(lessonWithTasks)
}