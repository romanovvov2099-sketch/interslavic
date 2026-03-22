package http

import (
	"interslavic/internal/auth"
	"interslavic/internal/http/handlers"
	"interslavic/internal/http/middlewares"
	usecase "interslavic/internal/usecases"
	"interslavic/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

// New /api/ Router
func NewRouter(app *fiber.App) fiber.Router {
	return app.Group("/api/")
}

// Register `/api` routes
func Register(
	router fiber.Router,
	logger *logging.ModuleLogger,
	authUseCase *usecase.AuthUseCase,
	courseUseCase *usecase.CourseUseCase,
	lessonUseCase *usecase.LessonUseCase,
	taskUseCase *usecase.TaskUseCase,
	progressUseCase *usecase.ProgressUseCase,
	jwtCfg *auth.JWTConfig,
) {
	mlogger := logging.NewModuleLogger("HTTP", "ROUTER", logger)
	mlogger.Info("register http api router")

	router.Use(middlewares.NewLoggingMiddleware(mlogger))
	router.Use(recover.New())
	router.Get("/healthy", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusOK) })

	// Auth routes
	authHandler := handlers.NewAuthHandler(authUseCase)
	authGroup := router.Group("/auth")
	{
		authGroup.Post("/login", authHandler.Login)
		authGroup.Post("/register", authHandler.Register)
		authGroup.Post("/refresh", authHandler.RefreshToken)

		// Protected auth routes
		protectedAuth := authGroup.Use(middlewares.JWTAuthMiddleware(jwtCfg))
		protectedAuth.Get("/me", authHandler.GetCurrentUser)
	}

	// Course routes
	courseHandler := handlers.NewCourseHandler(courseUseCase)
	courseGroup := router.Group("/courses")
	{
		courseGroup.Get("/", courseHandler.GetAllCourses)
		courseGroup.Get("/:id", courseHandler.GetCourseByID)
	}

	// Lesson routes
	lessonHandler := handlers.NewLessonHandler(lessonUseCase)

	// Lessons by course
	courseGroup.Get("/:id/lessons", lessonHandler.GetCourseLessons)

	// Lesson endpoints
	lessonGroup := router.Group("/lessons")
	{
		lessonGroup.Get("/:id", lessonHandler.GetLessonByID)
		lessonGroup.Get("/:id/full", lessonHandler.GetLessonWithTasks)
	}

	// Task routes
	taskHandler := handlers.NewTaskHandler(taskUseCase)

	// Tasks by lesson
	lessonGroup.Get("/:id/tasks", taskHandler.GetLessonTasks)

	// Task check (protected)
	taskGroup := router.Group("/tasks").Use(middlewares.JWTAuthMiddleware(jwtCfg))
	{
		taskGroup.Post("/check", taskHandler.CheckTaskAnswer)
	}

	// Progress routes (protected)
	progressHandler := handlers.NewProgressHandler(progressUseCase)
	progressGroup := router.Group("/progress").Use(middlewares.JWTAuthMiddleware(jwtCfg))
	{
		progressGroup.Post("/update", progressHandler.UpdateProgress)
		progressGroup.Get("/", progressHandler.GetUserProgress)
		progressGroup.Get("/lesson/:lesson_id", progressHandler.GetLessonProgress)

		// Новые маршруты для прогресса по курсам
		progressGroup.Get("/courses", progressHandler.GetAllCoursesProgress)
		progressGroup.Get("/course/:course_id", progressHandler.GetCourseProgress)
		progressGroup.Get("/course/:course_id/full", progressHandler.GetCourseWithProgress)
	}
}

// APIModule provides fiber.Router api and register controllers
var APIModule = fx.Module("http",
	fx.Provide(NewRouter),
	fx.Invoke(
		Register,
	),
)
