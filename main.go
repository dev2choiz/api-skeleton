package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/dev2choiz/api-skeleton/internal/config"
	"github.com/dev2choiz/api-skeleton/middleware"
	"github.com/dev2choiz/api-skeleton/pkg/business"
	"github.com/dev2choiz/api-skeleton/pkg/cache"
	"github.com/dev2choiz/api-skeleton/pkg/db"
	"github.com/dev2choiz/api-skeleton/pkg/env"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
	"github.com/dev2choiz/api-skeleton/pkg/logger"
	"github.com/dev2choiz/api-skeleton/pkg/repository"
	"github.com/dev2choiz/api-skeleton/server"

	_ "net/http/pprof"
)

func main() {
	ctx := context.Background()
	appEnv := env.GetString("APP_ENV", "production")
	logLevel := env.GetString("LOG_LEVEL", "info")

	logger.InitLogger("logs/api.log", logger.GetZapLogLevel(logLevel), appEnv == "development")
	logger := logger.Get(ctx)
	defer errapp.Check(ctx, logger.Sync())

	conf, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	db, err := db.New(conf.PostgresHost, conf.PostgresUser, conf.PostgresPassword, conf.PostgresDatabase, conf.PostgresPort, nil)
	if err != nil {
		logger.Fatal("failed to connect to prosgres", zap.Error(err))
	}

	re, err := cache.New(conf)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	repo := repository.New(db)
	bu := business.NewBusiness(repo, re, conf.JWTSecret)
	ser := server.New(bu)

	route := chi.NewRouter()
	route.Use(middleware.CorrelationIDMiddleware(), middleware.RecoverMiddleware())

	secure := route.Group(nil)
	secure.Use(middleware.AuthenticateMiddleware(bu), middleware.LogMiddleware())

	public := route.Group(nil)
	public.Use(middleware.LogMiddleware())
	server.ApplyRoutes(ser, public, secure)

	go func() {
		pprofRouter := chi.NewRouter()
		pprofRouter.Mount("/debug", chimiddleware.Profiler())

		addr := fmt.Sprintf(":%s", conf.PPROFPort)
		logger.Sugar().Infof("Pprof listen %s", addr)

		if err := http.ListenAndServe(addr, pprofRouter); err != nil {
			logger.Error("pprof", zap.String("port", addr), zap.Error(err))
		}
	}()

	addr := fmt.Sprintf(":%s", conf.APIPort)
	logger.Sugar().Infof("Listen %s", addr)
	if err := http.ListenAndServe(addr, route); err != nil {
		logger.Fatal("failed to server app", zap.String("port", addr), zap.Error(err))
	}
}
