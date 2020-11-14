package storage

import "time"

type User struct {
	ID        string     `db:"user_id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  []byte     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Url struct {
	URLHash     string    `db:"url_hash"`
	UserID      string    `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiredAt   time.Time `db:"expired_at"`
	OriginalURL string    `db:"original_url"`
	CustomAlias string    `db:"custom_alias"`
}
