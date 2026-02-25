package errapp

import (
	"context"

	"github.com/dev2choiz/api-skeleton/pkg/logger"
	"go.uber.org/zap"
)

func Check(ctx context.Context, err error) {
	if err != nil {
		logger.Get(ctx).Error("check", zap.Error(err))
	}
}
