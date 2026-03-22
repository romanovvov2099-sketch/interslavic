package models

import "time"

// Progress - модель прогресса обучения
type LessonProgress struct {
	ID             int        `json:"id" db:"id"`
	UserID         int        `json:"user_id" db:"user_id"`
	LessonID       int        `json:"lesson_id" db:"lesson_id"`
	Status         string     `json:"status" db:"status"`
	Score          int        `json:"score" db:"score"`
	CompletionDate *time.Time `json:"completion_date" db:"completion_date"`
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// UpdateProgressRequest - запрос на обновление прогресса
type UpdateProgressRequest struct {
	LessonID int    `json:"lesson_id" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=not_started in_progress completed"`
	Score    int    `json:"score" validate:"min=0,max=100"`
}
