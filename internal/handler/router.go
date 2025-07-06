package handler

import (
	"log/slog"
	"net/http"
	"time"
)

func NewRouter(orderHandler *OrderHandler, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("GET /order/{order_uid}/", orderHandler.GetOrder)

	return loggingMiddleware(mux, log)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(&l, r)

		log.Info("Got request", slog.Any("method", r.Method), slog.Any("path", r.URL), slog.Any("status_code", l.statusCode), slog.Any("time", time.Since(start)))
	})
}
