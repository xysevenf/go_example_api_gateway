package rate_limiter

import "sync"

const (
	baseLimit  = 8
	baseWindow = 60
)

type RateLimiter struct {
	limits limitsStorage
	mu     sync.Mutex
}

type limitsStorage interface {
	GetLimitCounter(key string) (int, bool)
	SetLimitCounter(key string, val int, ttl int) bool
}

type LimiterOptions struct {
	Storage limitsStorage
}

func NewLimiterOptions() *LimiterOptions {
	return &LimiterOptions{Storage: NewFreecacheLimitsStorage()}
}

func (r *RateLimiter) Allow(discriminator string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter, ok := r.limits.GetLimitCounter(discriminator)
	if ok {
		if counter > baseLimit {
			return false
		}
	}
	r.limits.SetLimitCounter(discriminator, counter+1, baseWindow)
	return true
}

func NewRateLimiter(options *LimiterOptions) *RateLimiter {
	return &RateLimiter{limits: options.Storage}
}
