package middleware

import (
	"log"
	"net/http"
	"time"
)

func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.Host+r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.Host+r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
