package repository

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/pkg/db"
)

func applyUserFilters(q *bun.SelectQuery, filters entity.UserFilters) {
	if filters.ID != nil {
		q.Where("id = ?", *filters.ID)
	}
	if filters.Username != nil {
		q.Where("username = ?", *filters.Username)
	}
	if filters.Firstname != nil {
		q.Where("firstname = ?", *filters.Firstname)
	}
	if filters.Lastname != nil {
		q.Where("lastname = ?", *filters.Lastname)
	}
	if filters.Limit != nil {
		q.Limit(*filters.Limit)
	}
}

// GetUser retrieves a user by its ID.
// It returns a not found error if no user matches the given ID.
func (r repository) GetUser(ctx context.Context, id string) (entity.User, error) {
	user := entity.User{}

	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)

	return user, db.WrapDBErr("get user", err)
}

// GetUsers retrieves all users matching the provided filters.
// It returns an empty slice if no users match.
func (r repository) GetUsers(ctx context.Context, filters entity.UserFilters) ([]entity.User, error) {
	users := []entity.User{}

	q := r.db.NewSelect().Model(&users)

	applyUserFilters(q, filters)

	return users, db.WrapDBErr("get users", q.Scan(ctx))
}

// GetOneUser retrieves a single user matching the provided filters.
// It returns an error if no user or more than one user is found.
func (r repository) GetOneUser(ctx context.Context, filters entity.UserFilters) (entity.User, error) {
	q := r.db.NewSelect()

	applyUserFilters(q, filters)

	return scanOne[entity.User](ctx, "get one user", q)
}

// InsertUser inserts a new user into the database.
// It returns the inserted user or an error if the operation fails.
func (r repository) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	_, err := r.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return entity.User{}, db.WrapDBErr("insert user", err)
	}

	return user, nil
}
