package handlers

import (
	"interslavic/internal/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CourseHandler struct {
	courseUseCase *usecases.CourseUseCase
}

func NewCourseHandler(courseUseCase *usecases.CourseUseCase) *CourseHandler {
	return &CourseHandler{
		courseUseCase: courseUseCase,
	}
}

// GetAllCourses godoc
// @Summary Get all courses
// @Description Get list of all available courses
// @Tags courses
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500 {object} map[string]string
// @Router /api/courses [get]
func (h *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	courses, err := h.courseUseCase.GetAllCourses(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(courses)
}

// GetCourseByID godoc
// @Summary Get course by ID
// @Description Get course details by ID
// @Tags courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid course id",
		})
	}

	course, err := h.courseUseCase.GetCourseByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(course)
}
