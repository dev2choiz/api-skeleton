package errapp

import (
	"errors"
	"fmt"
)

var (
	ErrAppNotFound     = errors.New("not found")
	ErrAppConflict     = errors.New("conflict")
	ErrAppUnauthorized = errors.New("unauthorized")
	ErrAppBadRequest   = errors.New("bad request")
	ErrAppInternal     = errors.New("internal error")
)

type ErrApp struct {
	Err error
}

func (e ErrApp) Error() string {
	return e.Err.Error()
}

func (e ErrApp) Unwrap() error {
	return e.Err
}

func wrapErr(appErr error, err error, msgs ...string) error {
	if len(msgs) > 0 && msgs[0] != "" {
		err = fmt.Errorf("%s: %w", msgs[0], err)
	}

	return ErrApp{
		Err: errors.Join(appErr, err),
	}
}

func WrapNotFound(err error, msg ...string) error {
	return wrapErr(ErrAppNotFound, err, msg...)
}

func WrapConflict(err error, msg ...string) error {
	return wrapErr(ErrAppConflict, err, msg...)
}

func WrapBadRequest(err error, msg ...string) error {
	return wrapErr(ErrAppBadRequest, err, msg...)
}

func WrapUnauthorized(err error, msg ...string) error {
	return wrapErr(ErrAppUnauthorized, err, msg...)
}

func WrapInternal(err error, msg ...string) error {
	return wrapErr(ErrAppInternal, err, msg...)
}

func newErr(appErr error, msg string, args ...any) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return ErrApp{
		Err: fmt.Errorf("%w %s", appErr, msg),
	}
}

func NewNotFound(msg string, args ...any) error {
	return newErr(ErrAppNotFound, msg, args...)
}

func NewConflict(msg string, args ...any) error {
	return newErr(ErrAppConflict, msg, args...)
}

func NewBadRequest(msg string, args ...any) error {
	return newErr(ErrAppBadRequest, msg, args...)
}

func NewUnauthorized(msg string, args ...any) error {
	return newErr(ErrAppUnauthorized, msg, args...)
}

func NewInternal(msg string, args ...any) error {
	return newErr(ErrAppInternal, msg, args...)
}
