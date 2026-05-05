package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pratikkumar2201/pk-proxy/util"
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
		fmt.Printf("Served Request %s \n", util.RequestURI(r))
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.Host+r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
