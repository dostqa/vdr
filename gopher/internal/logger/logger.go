package logger

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewLogger(env string) *slog.Logger {
	var logger *slog.Logger
	var options *slog.HandlerOptions

	switch env {
	case envLocal:
		options = &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					t := a.Value.Time()
					return slog.String(slog.TimeKey, t.Format("04:05"))
				}
				return a
			},
		}

		logger = slog.New(slog.NewTextHandler(os.Stdout, options))
	case envDev:
		options = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}

		logger = slog.New(slog.NewJSONHandler(os.Stdout, options))
	case envProd:
		options = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}

		logger = slog.New(slog.NewJSONHandler(os.Stdout, options))
	}

	return logger
}

func NewMiddlewareLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := logger.With(
			slog.String("component", "middleware_logger"),
		)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			defer func() {
				status := ww.Status()
				if status == 0 {
					status = http.StatusOK
				}

				args := []any{
					slog.Int("status", status),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("duration", time.Since(start)),
				}

				switch {
				case status >= 500:
					entry.Error("request completed", args...)
				case status >= 400:
					entry.Warn("request completed", args...)
				default:
					entry.Info("request completed", args...)
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
