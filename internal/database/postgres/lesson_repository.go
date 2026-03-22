package postgres

import (
	"context"
	"database/sql"
	"interslavic/internal/database"
	"interslavic/internal/models"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) database.LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) FindByID(ctx context.Context, id int) (*models.Lesson, error) {
	query := `SELECT id, course_id, title, content, multimedia, position 
	          FROM lessons 
	          WHERE id = $1`

	var lesson models.Lesson
	var position int
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Title,
		&lesson.Content,
		&lesson.Multimedia,
		&position,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	lesson.Position = position
	return &lesson, nil
}

func (r *LessonRepository) FindByCourseID(ctx context.Context, courseID int) ([]models.Lesson, error) {
	query := `SELECT id, course_id, title, content, multimedia, position 
	          FROM lessons 
	          WHERE course_id = $1 
	          ORDER BY position`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var lesson models.Lesson
		var position int
		err := rows.Scan(&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content, &lesson.Multimedia, &position)
		if err != nil {
			return nil, err
		}
		lesson.Position = position
		lessons = append(lessons, lesson)
	}

	return lessons, rows.Err()
}
