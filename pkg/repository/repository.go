package repository

import (
	"context"
	"errors"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/pkg/db"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
	"github.com/uptrace/bun"
)

type Repository interface {
	GetUser(ctx context.Context, id string) (entity.User, error)
	GetUsers(ctx context.Context, filters entity.UserFilters) ([]entity.User, error)
	GetOneUser(ctx context.Context, filters entity.UserFilters) (entity.User, error)
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
}

type repository struct {
	db *bun.DB
}

var (
	ErrNotFound       = errors.New("not found")
	ErrTooManuResults = errors.New("too many results")
)

func New(db *bun.DB) Repository {
	return &repository{db}
}

func scanOne[T any](ctx context.Context, action string, q *bun.SelectQuery) (T, error) {
	var zero T
	items := []T{}

	err := q.Model(&items).Scan(ctx)
	if err != nil {
		return zero, db.WrapDBErr(action, err)
	}

	if len(items) == 0 {
		return zero, errapp.WrapNotFound(ErrNotFound)
	} else if len(items) > 1 {
		return zero, errapp.WrapBadRequest(ErrTooManuResults)
	}

	return items[0], nil
}
