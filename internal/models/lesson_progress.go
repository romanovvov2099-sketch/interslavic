package models

import "time"

// Progress - модель прогресса обучения
type LessonProgress struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	LessonID       int       `json:"lesson_id" db:"lesson_id"`
	Status         string    `json:"status" db:"status"`
	Score          int       `json:"score" db:"score"`
	CompletionDate time.Time `json:"completion_date" db:"completion_date"`
}
