package repository

import (
	"os"
	"testing"

	"github.com/dev2choiz/api-skeleton/internal/testdb"
	"github.com/dev2choiz/api-skeleton/pkg/db"
	"github.com/dev2choiz/api-skeleton/pkg/fixtures"
	"github.com/uptrace/bun"
)

var dbTest *bun.DB

func newTestRepository() Repository {
	return New(dbTest)
}

func resetDB(t *testing.T) {
	t.Helper()

	ctx := t.Context()

	testdb.EnsureTestDatabase(t, ctx, dbTest)
	testdb.ResetSchema(t, ctx, dbTest)
	testdb.ApplyMigrations(t, ctx, dbTest)

	if _, err := dbTest.NewInsert().Model(&fixtures.Users).Exec(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	db, err := db.NewTest()
	if err != nil {
		panic(err)
	}

	dbTest = db

	code := m.Run()

	_ = db.Close()

	os.Exit(code)
}
