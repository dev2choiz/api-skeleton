package business

import (
	"context"

	"github.com/dev2choiz/api-skeleton/entity"
)

func (b *business) GetUsers(ctx context.Context, filters entity.UserFilters) ([]entity.User, error) {
	return b.repository.GetUsers(ctx, filters)
}
