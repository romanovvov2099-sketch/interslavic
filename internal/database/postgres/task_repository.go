package postgres

import (
	"context"
	"database/sql"
	"errors"
	"interslavic/internal/database"
	"interslavic/internal/models"
	"strings"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) database.TaskRepository {
	return &TaskRepository{db: db}
}

// parsePostgresArray конвертирует PostgreSQL массив в []string
// Формат: {value1,value2} или {"value with spaces","another value"}
func parsePostgresArray(data []byte) []string {
	if data == nil || len(data) == 0 {
		return []string{}
	}

	str := string(data)

	// Убираем фигурные скобки
	if len(str) < 2 || str[0] != '{' || str[len(str)-1] != '}' {
		return []string{}
	}

	str = str[1 : len(str)-1]

	if str == "" {
		return []string{}
	}

	var result []string
	var current strings.Builder
	inQuote := false
	escaped := false

	for i := 0; i < len(str); i++ {
		ch := str[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			escaped = true
		case '"':
			inQuote = !inQuote
		case ',':
			if !inQuote {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
		default:
			current.WriteByte(ch)
		}
	}

	// Добавляем последний элемент
	result = append(result, current.String())

	return result
}

func (r *TaskRepository) FindByID(ctx context.Context, id int) (*models.Task, error) {
	query := `SELECT id, lesson_id, task_type, question, answer, choises 
	          FROM tasks 
	          WHERE id = $1`

	var task models.Task
	var choisesBytes []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.LessonID,
		&task.TaskType,
		&task.Question,
		&task.Answer,
		&choisesBytes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Конвертируем []byte в []string
	task.Choises = parsePostgresArray(choisesBytes)

	return &task, nil
}

func (r *TaskRepository) FindByLessonID(ctx context.Context, lessonID int) ([]models.Task, error) {
	query := `SELECT id, lesson_id, task_type, question, answer, choises 
	          FROM tasks 
	          WHERE lesson_id = $1 
	          ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var choisesBytes []byte

		err := rows.Scan(
			&task.ID,
			&task.LessonID,
			&task.TaskType,
			&task.Question,
			&task.Answer,
			&choisesBytes,
		)
		if err != nil {
			return nil, err
		}

		// Конвертируем []byte в []string
		task.Choises = parsePostgresArray(choisesBytes)

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *TaskRepository) CheckAnswer(ctx context.Context, taskID int, answer string) (bool, error) {
	task, err := r.FindByID(ctx, taskID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, errors.New("task not found")
	}

	// Проверка ответа
	return answer == task.Answer, nil
}
