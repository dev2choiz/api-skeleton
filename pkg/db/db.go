package db

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/dev2choiz/api-skeleton/pkg/env"
	"github.com/dev2choiz/api-skeleton/pkg/errapp"
)

func New(host, user, password, database string, port int, tls *tls.Config) (*bun.DB, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", host, port)),
		pgdriver.WithTLSConfig(tls),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database),
	)

	sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	return db, db.Ping()
}

func NewTest() (*bun.DB, error) {
	port, err := env.GetInt("POSTGRES_PORT_TEST", 9433)
	if err != nil {
		return nil, err
	}

	return New(
		env.GetString("POSTGRES_HOST_TEST", "localhost"),
		env.GetString("POSTGRES_USER_TEST", "test"),
		env.GetString("POSTGRES_PASSWORD_TEST", "test"),
		env.GetString("POSTGRES_DATABASE_TEST", "test"),
		port,
		nil,
	)
}

func isIntegrityViolation(err error) bool {
	if pgDriverErr, ok := errors.AsType[pgdriver.Error](err); ok {
		return pgDriverErr.IntegrityViolation()
	}

	return false
}

func isTimeout(err error) bool {
	if pgDriverErr, ok := errors.AsType[pgdriver.Error](err); ok {
		return pgDriverErr.StatementTimeout()
	}

	return false
}

func isNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// WrapDBErr converts database errors into application domain errors.
func WrapDBErr(action string, err error) error {
	if err == nil {
		return nil
	}

	if isIntegrityViolation(err) {
		return errapp.WrapConflict(err, action)
	}

	if isTimeout(err) {
		return errapp.WrapInternal(err, action)
	}

	if isNotFound(err) {
		return errapp.WrapNotFound(err, action)
	}

	return errapp.WrapInternal(err, action)
}
