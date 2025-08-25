package limiter

import (
	"context"
	"time"
)

type FixedWindowLimiter struct {
	backend   FixedWindowBackend
	window    time.Duration
	limit     int64
	keyPrefix string
	clock     Clock
}

func NewFixedWindowLimiter(storage FixedWindowBackend, cfg FixedWindowConfig) RateLimiter {
	cfg = normalizeFixedWindowConfig(cfg)
	return &FixedWindowLimiter{backend: storage, window: cfg.Window, limit: cfg.Limit, keyPrefix: cfg.KeyPrefix, clock: realClock{}}
}

func (f *FixedWindowLimiter) Allow(ctx context.Context, key string) (Decision, error) {
	if key == "" {
		key = "unknown"
	}
	count, resetAt, err := f.backend.IncrWindow(ctx, f.keyPrefix+key, f.window, f.clock.Now())
	if err != nil {
		return Decision{Allowed: false}, err
	}
	if count > f.limit {
		return Decision{Allowed: false, Remaining: 0, ResetAt: resetAt, RetryAfter: time.Until(resetAt)}, nil
	}
	return Decision{Allowed: true, Remaining: f.limit - count, ResetAt: resetAt}, nil
}
