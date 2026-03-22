package models

// Lesson - модель урока
type Lesson struct {
	ID         int     `json:"id" db:"id"`
	CourseID   int     `json:"course_id" db:"course_id"`
	Title      string  `json:"title" db:"title"`
	Content    string  `json:"content" db:"content"`
	Multimedia *string `json:"multimedia" db:"multimedia"`
	Position   int     `json:"position" db:"position"`
}
