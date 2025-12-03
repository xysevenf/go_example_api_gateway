package rate_limiter

import (
	"testing"
)

func TestAllow_GoodCase(t *testing.T) {
	limiterService := NewRateLimiter(NewLimiterOptions())

	t.Run("genuine request", func(t *testing.T) {
		discriminator := "321"
		res := limiterService.Allow(discriminator)
		if !res {
			t.Errorf("good request dissallowed")
		}
	})
}

func TestAllow_BadCase(t *testing.T) {
	limiterOptions := NewLimiterOptions()
	limiterOptions.MaxReq = 1
	limiterService := NewRateLimiter(limiterOptions)

	t.Run("spam request", func(t *testing.T) {
		discriminator := "321"
		_ = limiterService.Allow(discriminator)
		res := limiterService.Allow(discriminator)
		if res {
			t.Errorf("repeated request allowed")
		}
	})
}
