package models

type TaskType int

const (
	ChoisesTask = 1
	WriteTask   = 2
)

type Task struct {
	ID       int      `json:"id" db:"id"`
	LessonID int      `json:"lesson_id" db:"lesson_id"`
	TaskType int      `json:"task_type" db:"task_type"`
	Question string   `json:"question" db:"question"`
	Answer   string   `json:"answer" db:"answer"`
	Choises  []string `json:"choises" db:"choises"`
}
