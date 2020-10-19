package storage

import (
	"context"
	"database/sql"
	"time"

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

	// ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("User not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidUserID = errors.New("ID is not in its proper form")
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
		if constraintError(dbErr) == ErrEmailAlreadyExist {
			return nil, ErrEmailAlreadyExist
		}

		return nil, errors.Wrapf(err, "inserting user")
	}

	return &u, nil
}

// RetrieveUser returns user with given user_id or ErrNotFound if no user was found by
// given user_id. Also can return ErrInvalidUserID if provided user_id is not uuid.
func RetrieveUser(ctx context.Context, db *sqlx.DB, user_id string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.RetrieveUser")
	defer span.End()

	if _, err := uuid.Parse(user_id); err != nil {
		return nil, ErrInvalidUserID
	}

	const q = `SELECT * FROM users where user_id = $1 and deleted_at is null`
	var u User

	if err := db.GetContext(ctx, &u, q, user_id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrapf(err, "selecting user with id %q", user_id)
	}

	return &u, nil
}

// DoesUserNameExist returns info about existing user name in database.
func DoesUserEmailExist(ctx context.Context, db *sqlx.DB, email string) (bool, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.DoesUserNameExist")
	defer span.End()

	var exist bool
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`

	err := db.GetContext(ctx, &exist, q, email)
	if err != nil {
		return exist, errors.Wrapf(err, "selecting user name exists %q", email)
	}

	return exist, err
}

// UpdateUserInfo replaces a user document in the database.
func UpdateUserInfo(ctx context.Context, db *sqlx.DB, userID, userName, email string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.Update")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `UPDATE users SET name = $2, email = $3 WHERE user_id = $1;`

	if _, err := db.ExecContext(ctx, q, userID, userName, email); err != nil {
		return constraintError(err)
	}

	return nil
}

// DeleteUser removes a user from the database.
func DeleteUser(ctx context.Context, db *sqlx.DB, userID string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.Delete")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `UPDATE users SET deleted_at=$1;`

	if _, err := db.ExecContext(ctx, q, time.Now()); err != nil {
		return errors.Wrapf(err, "deleting user %s", userID)
	}

	return nil
}

// UpdateUsersPassword updates users password with given userID.
func UpdateUsersPassword(ctx context.Context, db *sqlx.DB, userID, password string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.UpdateUsersPassword")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `UPDATE users SET password = $2 WHERE user_id = $1`

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}

	if _ , err := db.ExecContext(ctx, q, userID, hash); err != nil {
		return errors.Wrapf(err, "updating users password: %s", userID)
	}

	return nil
}

func constraintError(err error) error {
	const UniqueViolationCode = "23505"
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == UniqueViolationCode {
			switch pqErr.Constraint {
			case "users_email_key":
				return ErrEmailAlreadyExist
			}
		}
	}
	return err
}
