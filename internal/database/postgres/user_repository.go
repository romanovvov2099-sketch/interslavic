
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"interslavic/internal/models"
	datab "interslavic/internal/database"
	
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) datab.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (fullname, email, login, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, reg_date
	`

	err := r.db.QueryRowContext(ctx, query,
		user.Fullname,
		user.Email,
		user.Login,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.RegDate)

	return err
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT id, fullname, email, login, password, role, reg_date, last_login
		FROM users
		WHERE login = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Role,
		&user.RegDate,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, fullname, email, login, password, role, reg_date, last_login
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Role,
		&user.RegDate,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, fullname, email, login, password, role, reg_date, last_login
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Login,
		&user.Password,
		&user.Role,
		&user.RegDate,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int) error {
	query := `
		UPDATE users
		SET last_login = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}