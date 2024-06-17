package middleware

import (
	"net/http"
	"rate-limiter/limiter"
	"strings"
)

func RateLimitMiddleware(rl limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			token := r.Header.Get("API_KEY")
			allowed, err := rl.Allow(ip, token)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", rl.BlockDuration().String())
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
