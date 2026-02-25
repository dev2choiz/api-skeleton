package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
CREATE TABLE users (
  id uuid not null DEFAULT uuidv4() primary KEY,
	username text not null UNIQUE CHECK (username <> ''),
	password text not null CHECK (password <> ''),
	firstname text,
	lastname text,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz not null default current_timestamp
);`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`DROP TABLE users;`)
		return err
	})
}
