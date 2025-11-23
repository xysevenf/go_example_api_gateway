package rate_limiter

import "sync"

const (
	baseLimit  = 8
	baseWindow = 60
)

type RateLimiter struct {
	options *LimiterOptions
	mu      sync.Mutex
}

type limitsStorage interface {
	GetLimitCounter(key string) (int, bool)
	SetLimitCounter(key string, val int, ttl int) bool
}

type LimiterOptions struct {
	Storage limitsStorage
	MaxReq  int
	Window  int
}

func NewLimiterOptions() *LimiterOptions {
	return &LimiterOptions{Storage: NewFreecacheLimitsStorage(), MaxReq: baseLimit, Window: baseWindow}
}

func (r *RateLimiter) Allow(discriminator string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter, ok := r.options.Storage.GetLimitCounter(discriminator)
	if ok {
		if counter >= r.options.MaxReq {
			return false
		}
	}
	r.options.Storage.SetLimitCounter(discriminator, counter+1, r.options.Window)
	return true
}

func NewRateLimiter(options *LimiterOptions) *RateLimiter {
	return &RateLimiter{options: options}
}
