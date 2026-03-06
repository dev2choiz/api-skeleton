package main

import (
	"context"
	"os"

	"github.com/dev2choiz/api-skeleton/cmd/database"
	"github.com/dev2choiz/api-skeleton/internal/config"
	"github.com/dev2choiz/api-skeleton/pkg/db"
	"github.com/dev2choiz/api-skeleton/pkg/env"
	"github.com/dev2choiz/api-skeleton/pkg/logger"
	"go.uber.org/zap"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.Background()
	appEnv := env.GetString("APP_ENV", "production")
	logLevel := env.GetString("LOG_LEVEL", "info")

	logger.InitLogger("logs/cmd.log", logger.GetZapLogLevel(logLevel), appEnv == "development")
	logger := logger.Get(ctx)
	defer logger.Sync()

	conf, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	db, err := db.New(conf.PostgresHost, conf.PostgresUser, conf.PostgresPassword, conf.PostgresDatabase, conf.PostgresPort, nil)
	if err != nil {
		logger.Fatal("failed to connect to prosgres", zap.Error(err))
	}

	app := &cli.App{
		Name: "app",

		Commands: []*cli.Command{
			database.NewDatabaseCommand(context.Background(), db),
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Fatal("failed to run command", zap.Error(err))
	}
}
