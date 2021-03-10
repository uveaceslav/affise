package middleware

import (
	"net/http"
	"sync/atomic"
)

type RateLimitMiddleware interface {
	Limit(next http.HandlerFunc) http.HandlerFunc
}

type rateLimitMiddleware struct {
	rate    int32
	maxRate int
}

func NewRateLimitMiddleware(maxRate int) RateLimitMiddleware {
	return &rateLimitMiddleware{
		maxRate: maxRate,
	}
}

func (rl *rateLimitMiddleware) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&rl.rate) >= int32(rl.maxRate) {
			http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		atomic.AddInt32(&rl.rate, 1)
		next(rw, r)
		atomic.AddInt32(&rl.rate, -1)
	}
}
