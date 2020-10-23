package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Using constants in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Create users table.",
		Script: `
	CREATE TABLE IF NOT EXISTS users (
  		user_id UUID PRIMARY KEY, 
  		name VARCHAR(20) NOT NULL, 
  		email VARCHAR(255) NOT NULL UNIQUE,
  		password TEXT NOT NULL,
  		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  		updated_at TIMESTAMP DEFAULT NULL,
  		deleted_at TIMESTAMP DEFAULT NULL 
	);`,
	},
	{
		Version:     2,
		Description: "Create urls table.",
		Script: `
	CREATE TABLE IF NOT EXISTS urls (
  		url_hash VARCHAR(16) PRIMARY KEY,
  		user_id UUID NOT NULL,
  		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  		expired_at TIMESTAMP,
  		original_url TEXT NOT NULL,
  		custom_alias TEXT,
  	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id)
	);`,
	},
	{
		Version:     3,
		Description: "Create api_keys table",
		Script: `
	CREATE TABLE IF NOT EXISTS api_keys (
  		api_key UUID PRIMARY KEY,
  		user_id UUID NOT NULL,
  	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id)
	);`,
	},
	{
		Version: 4,
		Description: "Remove DeletedAt column from users",
		Script: `ALTER TABLE users DROP COLUMN deleted_at`,
	},
}
