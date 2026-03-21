package models

type Course struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Theory      string `json:"theory" db:"theory"`
}
 