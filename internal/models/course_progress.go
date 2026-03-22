package models

// AGREGAE DATA FROM DB

// CourseProgress - прогресс по курсу
type CourseProgress struct {
	CourseID          int     `json:"course_id"`
	CourseTitle       string  `json:"course_title"`
	TotalLessons      int     `json:"total_lessons"`
	CompletedLessons  int     `json:"completed_lessons"`
	InProgressLessons int     `json:"in_progress_lessons"`
	NotStartedLessons int     `json:"not_started_lessons"`
	ProgressPercent   float64 `json:"progress_percent"`
	Status            string  `json:"status"` // not_started, in_progress, completed
}

// LessonWithProgress - урок с информацией о прогрессе
type LessonWithProgress struct {
	Lesson
	ProgressStatus string  `json:"progress_status"`
	ProgressScore  int     `json:"progress_score"`
	CompletedAt    *string `json:"completed_at,omitempty"`
}

// CourseWithProgress - курс со списком уроков и прогрессом
type CourseWithProgress struct {
	Course
	Lessons          []LessonWithProgress `json:"lessons"`
	ProgressPercent  float64              `json:"progress_percent"`
	CompletedLessons int                  `json:"completed_lessons"`
	TotalLessons     int                  `json:"total_lessons"`
	Status           string               `json:"status"`
}
