package business

import (
	"context"
	"net/http"

	"github.com/dev2choiz/api-skeleton/entity"
	"github.com/dev2choiz/api-skeleton/pkg/cache"
	"github.com/dev2choiz/api-skeleton/pkg/repository"
)

type Business interface {
	GetUsers(ctx context.Context, filters entity.UserFilters) ([]entity.User, error)

	Register(ctx context.Context, username, password string) (entity.User, error)
	Authenticate(ctx context.Context, username, password string) (string, error)
	ValidateToken(r *http.Request) (entity.User, error)
}

type business struct {
	repository repository.Repository
	redis      cache.Cache
	jwtSecret  string
}

func NewBusiness(st repository.Repository, redis cache.Cache, jwtSecret string) Business {
	return &business{st, redis, jwtSecret}
}
