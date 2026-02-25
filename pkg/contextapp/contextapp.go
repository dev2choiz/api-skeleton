package contextapp

import (
	"context"

	"github.com/dev2choiz/api-skeleton/entity"
)

type ContextKey string

const (
	ContextKeyUser          ContextKey = "user"
	ContextKeyCorrelationID ContextKey = "correlation_id"
)

func GetCorrelationID(ctx context.Context) string {
	if v := ctx.Value(ContextKeyCorrelationID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func GetUser(ctx context.Context) entity.User {
	if v := ctx.Value(ContextKeyUser); v != nil {
		if u, ok := v.(entity.User); ok {
			return u
		}
	}

	return entity.User{}
}
