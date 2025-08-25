package limiter

import (
	"context"
	"errors"
	"time"
)

// fixed window rate limiter
type (
	Decision struct {
		Allowed      bool
		RetryAfter   time.Duration
		BlockedUntil time.Time
		Remaining    int64
		ResetAt      time.Time
	}

	RateLimiter interface {
		Allow(ctx context.Context, key string) (Decision, error)
	}
	FixedWindowConfig struct {
		Window    time.Duration
		Limit     int64
		KeyPrefix string
	}

	FixedWindowBackend interface {
		IncrWindow(ctx context.Context, baseKey string, window time.Duration, now time.Time) (count int64, resetAt time.Time, err error)
	}
)

var ErrUnexpectedScriptResult = errors.New("rateLimiter: unexpected backend result")

// normalize methods for configs
func normalizeFixedWindowConfig(cfg FixedWindowConfig) FixedWindowConfig {
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.Limit <= 0 {
		cfg.Limit = 60
	}
	return cfg
}
