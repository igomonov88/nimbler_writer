package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ErrURLNotFound is used when url is requested by the key does not exist.
var ErrURLNotFound = errors.New("Url not found")

// StoreURL used to store created alias with it's original url path and additional data.
func StoreURL(ctx context.Context, db *sqlx.DB, url Url) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.StoreURL")
	defer span.End()

	const query = `INSERT INTO urls (url_hash, user_id, expired_at, original_url, custom_alias) VALUES 
	(:url_hash, :user_id, :expired_at, :original_url, :custom_alias)`

	if _, err := db.NamedExecContext(ctx, query, url); err != nil {
		return errors.Wrapf(err, "inserting url with params: %v \n", url)
	}

	return nil
}

// RetrieveOriginalURL returned original urlPath if it matched with provided key, or
// ErrURLNotFound if there is no path for such key.
func RetrieveOriginalURL(ctx context.Context, db *sqlx.DB, key string) (string, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.RetrieveOriginalURL")
	defer span.End()

	const query = `SELECT original_url FROM urls WHERE url_hash = $1 AND expired_at < $2`
	var urlPath string

	if err := db.GetContext(ctx, &urlPath, query, key, time.Now()); err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", ErrURLNotFound
		default:
			return "",
				errors.Wrap(err, "selecting from urls")
		}
	}

	return urlPath, nil
}

// RetrieveAllExpiredURLKeysFromDate return all url hashes which expired.
func RetrieveAllExpiredURLKeysFromDate(ctx context.Context, db *sqlx.DB, date time.Time, limit int) ([]string, error) {
	ctx, span := trace.StartSpan(ctx, "internal.storage.RetrieveAllExpiredURLKeysFromDate")
	defer span.End()

	const query = `SELECT url_hash FROM urls WHERE expired_at < $1 limit $2`
	var urlHashes []string

	if err := db.SelectContext(ctx, &urlHashes, query, date, limit); err != nil {
		return nil, errors.Wrap(err, "selecting expired url hashes")
	}

	return urlHashes, nil
}

// DeleteURLS is deleting urls with provided hashes.
func DeleteURLS(ctx context.Context, db *sqlx.DB, hashes []string) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.DeleteURLS")
	defer span.End()

	const query = `DELETE FROM urls where url_hash = $1`

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin transaction.")
	}

	for i := range hashes {
		if _, err := tx.ExecContext(ctx, query, hashes[i]); err != nil {
			if err := tx.Rollback(); err != nil {
				return errors.Wrap(err, "rollback transaction.")
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit transaction result.")
	}

	return nil
}

// DeleteURL is deleting url with provided hash.
func DeleteURL(ctx context.Context, db *sqlx.DB, hash string) error {
	ctx, span := trace.StartSpan(ctx, "internal.storage.DeleteURL")
	defer span.End()

	const query = `DELETE FROM urls where url_hash = $1`

	if _, err := db.ExecContext(ctx, query, hash); err != nil {
		return errors.Wrapf(err, "delete url.")
	}

	return nil
}

// DoesURLExist returns info about existing url in database.
func DoesURLAliasExist(ctx context.Context, db *sqlx.DB, alias string) (bool, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.DoesURLAliasExist")
	defer span.End()

	var exist bool
	const q = `SELECT EXISTS(SELECT 1 FROM urls WHERE custom_alias = $1);`

	err := db.GetContext(ctx, &exist, q, alias)
	if err != nil {
		return exist, errors.Wrapf(err, "selecting custom alias exist %q", alias)
	}

	return exist, err
}
