package server

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/dev2choiz/api-skeleton/pkg/logger"
)

func getQueryString(r *http.Request, key string) *string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}

	return &v
}

func getQueryInt(r *http.Request, key string) *int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		logger.Get(r.Context()).Error("error while parsing query to int", zap.String("key", key), zap.String("value", v))

		return new(0)
	}

	return &n
}
