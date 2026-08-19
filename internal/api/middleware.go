package api

import (
	"log/slog"
	"net/http"
	"time"
)

// connect routers to the sustom logger
func withLogging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w , code: http.StatusOK} // default to 200 OK
			next.ServeHTTP(rec, r)
			logger.Info("api request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.code,
				"latency_ms", time.Since(start).Milliseconds(),
		
		)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}