package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	"github.com/dev2choiz/api-skeleton/cmd/database/migrations"
	"github.com/dev2choiz/api-skeleton/pkg/fixtures"
	"github.com/dev2choiz/api-skeleton/pkg/logger"
)

func NewDatabaseCommand(ctx context.Context, db *bun.DB) *cli.Command {
	logger := logger.Get(ctx)
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	return &cli.Command{
		Name:  "database",
		Usage: "database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "fixtures",
				Usage: "run fixtures",
				Action: func(c *cli.Context) error {
					_, err := db.NewInsert().Model(&fixtures.Users).Exec(ctx)

					return err
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context) //nolint:errcheck

					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						logger.Info("there are no new migrations to run (database is up to date)")

						return nil
					}
					logger.Sugar().Infof("migrated to %s", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context) //nolint:errcheck

					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						logger.Info("there are no groups to roll back")
						return nil
					}
					logger.Sugar().Infof("rolled back %s", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					return migrator.Lock(c.Context)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					return migrator.Unlock(c.Context)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(c.Context, name)
					if err != nil {
						return err
					}
					logger.Sugar().Infof("created migration %s (%s)", mf.Name, mf.Path)
					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						logger.Sugar().Infof("created migration %s (%s)", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "create_tx_sql",
				Usage: "create up and down transactional SQL migrations",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateTxSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						logger.Sugar().Infof("created transaction migration %s (%s)", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ms, err := migrator.MigrationsWithStatus(c.Context)
					if err != nil {
						return err
					}
					logger.Sugar().Infof("migrations: %s", ms)
					logger.Sugar().Infof("unapplied migrations: %s", ms.Unapplied())
					logger.Sugar().Infof("last migration group: %s", ms.LastGroup())
					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					group, err := migrator.Migrate(c.Context, migrate.WithNopMigration())
					if err != nil {
						return err
					}
					if group.IsZero() {
						logger.Info("there are no new migrations to mark as applied")
						return nil
					}
					logger.Sugar().Infof("marked as applied %s", group)
					return nil
				},
			},
			{
				Name:  "wait",
				Usage: "wait for database to be ready",
				Flags: []cli.Flag{
					&cli.DurationFlag{
						Name:    "timeout",
						Aliases: []string{"t"},
						Value:   30 * time.Second,
						Usage:   "maximum time to wait for database",
					},
				},
				Action: func(c *cli.Context) error {
					timeout := c.Duration("timeout")
					return waitDatabase(c.Context, db, timeout)
				},
			},
		},
	}
}

func waitDatabase(ctx context.Context, db *bun.DB, timeout time.Duration) error {
	logger := logger.Get(ctx)
	deadline := time.Now().Add(timeout)

	for {
		if err := db.PingContext(ctx); err == nil {
			logger.Info("Database is ready!")
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timeout reached: database is not ready after %s", timeout)
		}

		time.Sleep(1 * time.Second)
	}
}
