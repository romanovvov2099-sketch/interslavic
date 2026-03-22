package postgres

import (
	"context"
	"database/sql"
	"interslavic/internal/database"
	"interslavic/internal/models"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) database.CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) FindAll(ctx context.Context) ([]models.Course, error) {
	query := `SELECT id, title, description, theory FROM courses ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.Theory)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, rows.Err()
}

func (r *CourseRepository) FindByID(ctx context.Context, id int) (*models.Course, error) {
	query := `SELECT id, title, description, theory FROM courses WHERE id = $1`

	var course models.Course
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Theory,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &course, nil
}
