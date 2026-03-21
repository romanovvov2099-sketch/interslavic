package models

import "time"

// User - модель пользователя
type User struct {
	ID        int       `json:"id" db:"id"`
	Fullname  string    `json:"fullname" db:"fullname"`
	Email     string    `json:"email" db:"email"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"-" db:"password"` // "-" скрывает поле при JSON сериализации
	Role      string    `json:"role" db:"role"`
	RegDate   time.Time `json:"reg_date" db:"reg_date"`
	LastLogin time.Time `json:"last_login" db:"last_login"`
}
