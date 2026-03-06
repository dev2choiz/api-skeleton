package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/dev2choiz/api-skeleton/pkg/business"
	"github.com/dev2choiz/api-skeleton/pkg/contextapp"
	"github.com/dev2choiz/api-skeleton/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LogMiddleware logs HTTP requests including method, path, status code,
// duration, client IP, and user agent.
func LogMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			logger.Get(r.Context()).Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.status),
				zap.Duration("duration", duration),
				zap.String("ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}

// AuthenticateMiddleware validates the JWT token and injects the user
// into the request context. It returns 401 if authentication fails.
func AuthenticateMiddleware(bu business.Business) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := bu.ValidateToken(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), contextapp.ContextKeyUser, user))

			next.ServeHTTP(w, r)
		})
	}
}

// RecoverMiddleware recovers from panics, logs the stack trace,
// and returns a 500 Internal Server Error response.
func RecoverMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Get(r.Context()).Error(
						"panic recovered",
						zap.Any("panic", rec),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("ip", r.RemoteAddr),
						zap.String("user_agent", r.UserAgent()),
						zap.Stack("stack"),
					)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func CorrelationIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cid := r.Header.Get("X-Correlation-ID")
			if cid == "" {
				cid = uuid.NewString()
			}

			ctx := context.WithValue(r.Context(), contextapp.ContextKeyCorrelationID, cid)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
