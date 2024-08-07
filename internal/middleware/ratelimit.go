package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// only allows rps number of requests every per time(minute, second)
func RateLimit(rps int, per time.Duration) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(per/time.Duration(rps)), rps)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
