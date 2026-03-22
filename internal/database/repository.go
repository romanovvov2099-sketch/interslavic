package database

import (
	"context"
	"database/sql"
	"interslavic/internal/models"
)

// [SERVICE INTERFACES]

type TrManager interface {
	Begin(context.Context) (context.Context, error)
	Commit(context.Context) error
	Rollback(context.Context) error
	GetCurrentTr(context.Context) *sql.Tx
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID int) error
}

type CourseRepository interface {
	FindAll(ctx context.Context) ([]models.Course, error)
	FindByID(ctx context.Context, id int) (*models.Course, error)
}

type LessonRepository interface {
	FindByID(ctx context.Context, id int) (*models.Lesson, error)
	FindByCourseID(ctx context.Context, courseID int) ([]models.Lesson, error)
}

type TaskRepository interface {
	FindByID(ctx context.Context, id int) (*models.Task, error)
	FindByLessonID(ctx context.Context, lessonID int) ([]models.Task, error)
	CheckAnswer(ctx context.Context, taskID int, answer string) (bool, error)
}

type LessonProgressRepository interface {
	CreateOrUpdate(ctx context.Context, progress *models.LessonProgress) error
	FindByUserAndLesson(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error)
	FindByUserID(ctx context.Context, userID int) ([]models.LessonProgress, error)
	FindByUserAndCourse(ctx context.Context, userID, courseID int) ([]models.LessonProgress, error)
	GetCourseProgressStats(ctx context.Context, userID, courseID int) (*models.CourseProgress, error)
	GetAllCoursesProgress(ctx context.Context, userID int) ([]models.CourseProgress, error)
	GetLessonsWithProgressByCourse(ctx context.Context, userID, courseID int) ([]models.LessonWithProgress, error)
}
