package limiter

import (
	"time"
)

type RateLimiter interface {
	Allow(ip string, token string) (bool, error)
	BlockDuration() time.Duration
}
