package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestSlogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		elapsed := time.Since(start)
		slog.Info(
			"Incoming HTTP request:",
			slog.String("method", r.Method),
			slog.String("URI", r.RequestURI),
			slog.Int("statusCode", ww.Status()),
			slog.String("elapsedTime", elapsed.String()),
			slog.Int("responseSize", ww.BytesWritten()),
		)
	}
	return http.HandlerFunc(fn)
}
