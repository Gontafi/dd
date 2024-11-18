package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(
			"[%s] %s %s %s %d %v",
			r.Method,
			r.RemoteAddr,
			r.URL.Path,
			r.Proto,
			rw.status,
			duration,
		)
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
