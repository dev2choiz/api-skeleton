package testdb

import (
	"context"
	"strings"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	"github.com/dev2choiz/api-skeleton/cmd/database/migrations"
)

// EnsureTestDatabase prevents accidental execution on production database.
func EnsureTestDatabase(t *testing.T, ctx context.Context, db *bun.DB) {
	t.Helper()

	var dbName string
	if err := db.QueryRowContext(ctx, `SELECT current_database()`).Scan(&dbName); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(strings.ToLower(dbName), "test") {
		t.Fatalf("refusing to drop schema on non-test database: %s", dbName)
	}
}

// ResetSchema drops and recreates the public schema
func ResetSchema(t *testing.T, ctx context.Context, db *bun.DB) {
	t.Helper()

	if _, err := db.ExecContext(ctx, `DROP SCHEMA public CASCADE;`); err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `CREATE SCHEMA public;`); err != nil {
		t.Fatal(err)
	}
}

// ApplyMigrations runs all migrations
func ApplyMigrations(t *testing.T, ctx context.Context, db *bun.DB) {
	t.Helper()

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	if err := migrator.Init(ctx); err != nil {
		t.Fatal(err)
	}

	if err := migrator.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	defer migrator.Unlock(ctx)

	if _, err := migrator.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
}
