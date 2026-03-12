package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (ww *wrappedResponseWriter) WriteHeader(statusCode int) {
	ww.statusCode = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *wrappedResponseWriter) Write(data []byte) (int, error) {
	n, err := ww.ResponseWriter.Write(data)
	ww.bytesWritten += n
	return n, err
}

func Logger(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := newWrappedResponseWriter(w)

			start := time.Now()
			next.ServeHTTP(ww, r)

			log.Info(
				"request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.statusCode),
				zap.Int("bytes", ww.bytesWritten),
				zap.Duration("duration", time.Since(start)),
				// zap.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	}
}
