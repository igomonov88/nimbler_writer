package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailAlreadyExist is used when User with such email is already exist.
	ErrEmailAlreadyExist = errors.New("Email already exist")
)

// CreateUser creates user with provided info in database.
func CreateUser(ctx context.Context, db *sqlx.DB, name, email, password string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.CreateUser")
	defer span.End()

	const query = `INSERT INTO users (user_id, name, email, password) 
	VALUES (:user_id, :name, :email, :password)`

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}
	u := User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: hash,
	}
	_, dbErr := db.NamedExec(query, &u)
	if dbErr != nil {
		return nil, errors.Wrap(constraintError(err), "inserting user")
	}

	return &u, nil
}

func constraintError(err error) error {
	const UniqueViolationCode = "23505"
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == UniqueViolationCode {
			switch pqErr.Constraint {
			case "email_idx":
				return ErrEmailAlreadyExist
			}
		}
	}
	return err
}
