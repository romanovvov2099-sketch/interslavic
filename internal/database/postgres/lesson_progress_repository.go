package postgres

import (
	"context"
	"database/sql"
	"errors"
	"interslavic/internal/database"
	"interslavic/internal/models"
)

type LessonProgressRepository struct {
	db *sql.DB
}

func NewLessonProgressRepository(db *sql.DB) database.LessonProgressRepository {
	return &LessonProgressRepository{db: db}
}

func (r *LessonProgressRepository) CreateOrUpdate(ctx context.Context, progress *models.LessonProgress) error {
	// Проверяем, существует ли запись
	existing, err := r.FindByUserAndLesson(ctx, progress.UserID, progress.LessonID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Если запись существует - обновляем
	if existing != nil {
		query := `
			UPDATE lesson_progress
			SET status = $1, score = $2, completion_date = $3, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = $4 AND lesson_id = $5
			RETURNING id, started_at, updated_at
		`

		var completionDate sql.NullTime
		if progress.CompletionDate != nil {
			completionDate = sql.NullTime{Time: *progress.CompletionDate, Valid: true}
		}

		err = r.db.QueryRowContext(ctx, query,
			progress.Status,
			progress.Score,
			completionDate,
			progress.UserID,
			progress.LessonID,
		).Scan(&progress.ID, &progress.StartedAt, &progress.UpdatedAt)

		return err
	}

	// Если записи нет - создаем новую
	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, status, score, completion_date, started_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, started_at, updated_at
	`

	var completionDate sql.NullTime
	if progress.CompletionDate != nil {
		completionDate = sql.NullTime{Time: *progress.CompletionDate, Valid: true}
	}

	err = r.db.QueryRowContext(ctx, query,
		progress.UserID,
		progress.LessonID,
		progress.Status,
		progress.Score,
		completionDate,
	).Scan(&progress.ID, &progress.StartedAt, &progress.UpdatedAt)

	return err
}

func (r *LessonProgressRepository) FindByUserAndLesson(ctx context.Context, userID, lessonID int) (*models.LessonProgress, error) {
	query := `
		SELECT id, user_id, lesson_id, status, score, completion_date, started_at, updated_at
		FROM lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`

	var progress models.LessonProgress
	var completionDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, lessonID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.LessonID,
		&progress.Status,
		&progress.Score,
		&completionDate,
		&progress.StartedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if completionDate.Valid {
		progress.CompletionDate = &completionDate.Time
	}

	return &progress, nil
}

func (r *LessonProgressRepository) FindByUserID(ctx context.Context, userID int) ([]models.LessonProgress, error) {
	query := `
		SELECT id, user_id, lesson_id, status, score, completion_date, started_at, updated_at
		FROM lesson_progress
		WHERE user_id = $1
		ORDER BY lesson_id
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []models.LessonProgress
	for rows.Next() {
		var progress models.LessonProgress
		var completionDate sql.NullTime

		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.LessonID,
			&progress.Status,
			&progress.Score,
			&completionDate,
			&progress.StartedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if completionDate.Valid {
			progress.CompletionDate = &completionDate.Time
		}

		progresses = append(progresses, progress)
	}

	return progresses, rows.Err()
}


// FindByUserAndCourse - получить прогресс по всем урокам курса для пользователя
func (r *LessonProgressRepository) FindByUserAndCourse(ctx context.Context, userID, courseID int) ([]models.LessonProgress, error) {
	query := `
		SELECT lp.id, lp.user_id, lp.lesson_id, lp.status, lp.score, 
		       lp.completion_date, lp.started_at, lp.updated_at
		FROM lesson_progress lp
		INNER JOIN lessons l ON l.id = lp.lesson_id
		WHERE lp.user_id = $1 AND l.course_id = $2
		ORDER BY l.position
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var progresses []models.LessonProgress
	for rows.Next() {
		var progress models.LessonProgress
		var completionDate sql.NullTime
		
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.LessonID,
			&progress.Status,
			&progress.Score,
			&completionDate,
			&progress.StartedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if completionDate.Valid {
			progress.CompletionDate = &completionDate.Time
		}
		
		progresses = append(progresses, progress)
	}
	
	return progresses, rows.Err()
}

// GetCourseProgressStats - получить статистику прогресса по курсу
func (r *LessonProgressRepository) GetCourseProgressStats(ctx context.Context, userID, courseID int) (*models.CourseProgress, error) {
	query := `
		WITH course_lessons AS (
			SELECT id, position
			FROM lessons
			WHERE course_id = $1
		),
		user_progress AS (
			SELECT 
				lesson_id,
				status,
				score
			FROM lesson_progress
			WHERE user_id = $2
		)
		SELECT 
			c.id as course_id,
			c.title as course_title,
			COUNT(cl.id) as total_lessons,
			COUNT(CASE WHEN up.status = 'completed' THEN 1 END) as completed_lessons,
			COUNT(CASE WHEN up.status = 'in_progress' THEN 1 END) as in_progress_lessons,
			COUNT(CASE WHEN up.status = 'not_started' OR up.status IS NULL THEN 1 END) as not_started_lessons,
			ROUND(
				COUNT(CASE WHEN up.status = 'completed' THEN 1 END)::numeric / 
				NULLIF(COUNT(cl.id), 0) * 100, 
				2
			) as progress_percent
		FROM courses c
		CROSS JOIN course_lessons cl
		LEFT JOIN user_progress up ON up.lesson_id = cl.id
		WHERE c.id = $1
		GROUP BY c.id, c.title
	`
	
	var stats models.CourseProgress
	err := r.db.QueryRowContext(ctx, query, courseID, userID).Scan(
		&stats.CourseID,
		&stats.CourseTitle,
		&stats.TotalLessons,
		&stats.CompletedLessons,
		&stats.InProgressLessons,
		&stats.NotStartedLessons,
		&stats.ProgressPercent,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Определяем статус курса
	if stats.CompletedLessons == stats.TotalLessons {
		stats.Status = "completed"
	} else if stats.CompletedLessons > 0 || stats.InProgressLessons > 0 {
		stats.Status = "in_progress"
	} else {
		stats.Status = "not_started"
	}
	
	return &stats, nil
}

// GetAllCoursesProgress - получить прогресс по всем курсам для пользователя
func (r *LessonProgressRepository) GetAllCoursesProgress(ctx context.Context, userID int) ([]models.CourseProgress, error) {
	query := `
		WITH all_courses AS (
			SELECT 
				c.id as course_id,
				c.title as course_title,
				COUNT(l.id) as total_lessons
			FROM courses c
			LEFT JOIN lessons l ON l.course_id = c.id
			GROUP BY c.id, c.title
		),
		user_progress AS (
			SELECT 
				l.course_id,
				lp.status,
				COUNT(lp.id) as progress_count
			FROM lesson_progress lp
			INNER JOIN lessons l ON l.id = lp.lesson_id
			WHERE lp.user_id = $1
			GROUP BY l.course_id, lp.status
		)
		SELECT 
			ac.course_id,
			ac.course_title,
			ac.total_lessons,
			COALESCE(
				(SELECT progress_count FROM user_progress up 
				 WHERE up.course_id = ac.course_id AND up.status = 'completed'), 0
			) as completed_lessons,
			COALESCE(
				(SELECT progress_count FROM user_progress up 
				 WHERE up.course_id = ac.course_id AND up.status = 'in_progress'), 0
			) as in_progress_lessons,
			ac.total_lessons - 
				COALESCE((SELECT progress_count FROM user_progress up 
				 WHERE up.course_id = ac.course_id AND up.status IN ('completed', 'in_progress')), 0) 
			as not_started_lessons,
			ROUND(
				COALESCE(
					(SELECT progress_count FROM user_progress up 
					 WHERE up.course_id = ac.course_id AND up.status = 'completed'), 0
				)::numeric / NULLIF(ac.total_lessons, 0) * 100, 
				2
			) as progress_percent
		FROM all_courses ac
		ORDER BY ac.course_id
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var stats []models.CourseProgress
	for rows.Next() {
		var stat models.CourseProgress
		err := rows.Scan(
			&stat.CourseID,
			&stat.CourseTitle,
			&stat.TotalLessons,
			&stat.CompletedLessons,
			&stat.InProgressLessons,
			&stat.NotStartedLessons,
			&stat.ProgressPercent,
		)
		if err != nil {
			return nil, err
		}
		
		// Определяем статус курса
		if stat.CompletedLessons == stat.TotalLessons {
			stat.Status = "completed"
		} else if stat.CompletedLessons > 0 || stat.InProgressLessons > 0 {
			stat.Status = "in_progress"
		} else {
			stat.Status = "not_started"
		}
		
		stats = append(stats, stat)
	}
	
	return stats, rows.Err()
}

// GetLessonsWithProgressByCourse - получить уроки курса с прогрессом пользователя
func (r *LessonProgressRepository) GetLessonsWithProgressByCourse(ctx context.Context, userID, courseID int) ([]models.LessonWithProgress, error) {
	query := `
		SELECT 
			l.id,
			l.course_id,
			l.title,
			l.content,
			l.multimedia,
			l.position,
			COALESCE(lp.status, 'not_started') as progress_status,
			COALESCE(lp.score, 0) as progress_score,
			CASE 
				WHEN lp.completion_date IS NOT NULL THEN TO_CHAR(lp.completion_date, 'YYYY-MM-DD HH24:MI:SS')
				ELSE NULL 
			END as completed_at
		FROM lessons l
		LEFT JOIN lesson_progress lp ON lp.lesson_id = l.id AND lp.user_id = $1
		WHERE l.course_id = $2
		ORDER BY l.position
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var lessons []models.LessonWithProgress
	for rows.Next() {
		var lesson models.LessonWithProgress
		var completedAt *string
		
		err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Title,
			&lesson.Content,
			&lesson.Multimedia,
			&lesson.Position,
			&lesson.ProgressStatus,
			&lesson.ProgressScore,
			&completedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if completedAt != nil {
			lesson.CompletedAt = completedAt
		}
		
		lessons = append(lessons, lesson)
	}
	
	return lessons, rows.Err()
}