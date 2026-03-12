package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/internal/testdb"
	"github.com/dev2choiz/api-skeleton/pkg/db"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
	"github.com/dev2choiz/api-skeleton/pkg/fixtures"
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

func Test_repository_scanOne(t *testing.T) {
	ctx := context.Background()
	resetDB(t)

	tests := []struct {
		name       string
		setupQuery func() *bun.SelectQuery
		want       entity.User
		wantErr    error
	}{
		{
			name: "Error in the query",
			setupQuery: func() *bun.SelectQuery {
				q := dbTest.NewSelect().Model(&entity.User{})
				applyUserFilters(q, entity.UserFilters{ID: new("invaliduuid")})
				return q
			},
			want:    entity.User{},
			wantErr: errapp.ErrAppInternal,
		},
		{
			name: "Not found",
			setupQuery: func() *bun.SelectQuery {
				q := dbTest.NewSelect().Model(&entity.User{})
				applyUserFilters(q, entity.UserFilters{Username: new("doesnotexist")})
				return q
			},
			want:    entity.User{},
			wantErr: errapp.ErrAppNotFound,
		},
		{
			name: "Too many results",
			setupQuery: func() *bun.SelectQuery {
				q := dbTest.NewSelect().Model(&entity.User{})
				applyUserFilters(q, entity.UserFilters{})
				return q
			},
			want:    entity.User{},
			wantErr: errapp.ErrAppBadRequest,
		},
		{
			name: "Success",
			setupQuery: func() *bun.SelectQuery {
				q := dbTest.NewSelect().Model(&entity.User{})
				applyUserFilters(q, entity.UserFilters{Username: &fixtures.Users[0].Username})
				return q
			},
			want: fixtures.Users[0],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.setupQuery()
			got, err := scanOne[entity.User](ctx, "test", q)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, got.ID, tt.want.ID)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
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
