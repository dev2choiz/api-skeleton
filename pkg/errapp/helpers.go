package errapp

import (
	"context"

	"go.uber.org/zap"

	"github.com/dev2choiz/api-skeleton/pkg/logger"
)

func Check(ctx context.Context, err error) {
	if err != nil {
		logger.Get(ctx).Error("check", zap.Error(err))
	}
}
