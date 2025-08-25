package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func getRedis() *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "yourpassword",
		DB:       9,
	})
	return rds
}

// helper to build a limiter with a controllable clock
func newTestLimiter(limit int64, window time.Duration, prefix string, start time.Time) (*FixedWindowLimiter, *mockClock) {
	storage := NewMemoryStorage()
	cfg := FixedWindowConfig{Window: window, Limit: limit, KeyPrefix: prefix}
	l := NewFixedWindowLimiter(storage, cfg).(*FixedWindowLimiter)
	mc := NewMockClock(start)
	l.clock = mc
	return l, mc
}

// --- Redis helpers and tests ---

func checkRedisOrSkip(t *testing.T) *redis.Client {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cli := getRedis()
	if err := cli.Ping(ctx).Err(); err != nil {
		t.Skipf("skipping Redis-backed tests: cannot connect to redis: %v", err)
	}
	if err := cli.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to flush redis db: %v", err)
	}
	return cli
}

func newRedisLimiter(cli *redis.Client, limit int64, window time.Duration, prefix string, start time.Time) (*FixedWindowLimiter, *mockClock) {
	storage := NewRedisStorage(cli)
	cfg := FixedWindowConfig{Window: window, Limit: limit, KeyPrefix: prefix}
	l := NewFixedWindowLimiter(storage, cfg).(*FixedWindowLimiter)
	mc := NewMockClock(start)
	l.clock = mc
	return l, mc
}

func TestFixedWindowRedis_AllowWithinLimitAndBlock(t *testing.T) {
	cli := checkRedisOrSkip(t)
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(3)
	window := 10 * time.Second
	l, mc := newRedisLimiter(cli, limit, window, "r1:", start)

	ctx := context.Background()
	for i := int64(1); i <= limit; i++ {
		dec, err := l.Allow(ctx, "k")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !dec.Allowed {
			t.Fatalf("call %d expected allowed=true", i)
		}
		expectedRemaining := limit - i
		if dec.Remaining != expectedRemaining {
			t.Fatalf("call %d expected remaining=%d, got %d", i, expectedRemaining, dec.Remaining)
		}
		windowStart := mc.Now().Truncate(window)
		expectedReset := windowStart.Add(window)
		if !dec.ResetAt.Equal(expectedReset) {
			t.Fatalf("call %d expected resetAt=%v, got %v", i, expectedReset, dec.ResetAt)
		}
	}
	dec, err := l.Allow(ctx, "k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Allowed {
		t.Fatalf("expected allowed=false after limit exceeded")
	}
	if dec.Remaining != 0 {
		t.Fatalf("expected remaining=0 when blocked, got %d", dec.Remaining)
	}
	windowStart := mc.Now().Truncate(window)
	expectedReset := windowStart.Add(window)
	if !dec.ResetAt.Equal(expectedReset) {
		t.Fatalf("expected resetAt=%v, got %v", expectedReset, dec.ResetAt)
	}
}

func TestFixedWindowRedis_ResetsNextWindow(t *testing.T) {
	cli := checkRedisOrSkip(t)
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(2)
	window := 5 * time.Second
	l, mc := newRedisLimiter(cli, limit, window, "r2:", start)
	ctx := context.Background()

	for i := int64(0); i < limit; i++ {
		if dec, err := l.Allow(ctx, "user"); err != nil || !dec.Allowed {
			t.Fatalf("pre-window calls should be allowed: dec=%+v err=%v", dec, err)
		}
	}
	if dec, _ := l.Allow(ctx, "user"); dec.Allowed {
		t.Fatalf("expected blocked after exhausting limit")
	}

	mc.Advance(window)
	dec, err := l.Allow(ctx, "user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dec.Allowed {
		t.Fatalf("expected allowed at the start of new window")
	}
	if dec.Remaining != limit-1 {
		t.Fatalf("expected remaining=%d, got %d", limit-1, dec.Remaining)
	}
	expectedReset := mc.Now().Truncate(window).Add(window)
	if !dec.ResetAt.Equal(expectedReset) {
		t.Fatalf("expected resetAt=%v, got %v", expectedReset, dec.ResetAt)
	}
}

func TestFixedWindowRedis_KeyPrefixIsolation(t *testing.T) {
	cli := checkRedisOrSkip(t)
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(3)
	window := 10 * time.Second

	storage := NewRedisStorage(cli)
	cfgA := FixedWindowConfig{Window: window, Limit: limit, KeyPrefix: "a:"}
	cfgB := FixedWindowConfig{Window: window, Limit: limit, KeyPrefix: "b:"}
	la := NewFixedWindowLimiter(storage, cfgA).(*FixedWindowLimiter)
	lb := NewFixedWindowLimiter(storage, cfgB).(*FixedWindowLimiter)
	mc := NewMockClock(start)
	la.clock = mc
	lb.clock = mc

	ctx := context.Background()
	for i := int64(0); i < limit; i++ {
		if dec, err := la.Allow(ctx, "k"); err != nil || !dec.Allowed {
			t.Fatalf("expected la allowed, got dec=%+v err=%v", dec, err)
		}
	}
	if dec, _ := la.Allow(ctx, "k"); dec.Allowed {
		t.Fatalf("expected la blocked after limit")
	}
	dec, err := lb.Allow(ctx, "k")
	if err != nil || !dec.Allowed || dec.Remaining != limit-1 {
		t.Fatalf("expected lb allowed and independent, dec=%+v err=%v", dec, err)
	}
}

func TestFixedWindowRedis_EmptyKeyUsesUnknownBucket(t *testing.T) {
	cli := checkRedisOrSkip(t)
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(2)
	window := 10 * time.Second
	l, _ := newRedisLimiter(cli, limit, window, "r4:", start)
	ctx := context.Background()

	dec, err := l.Allow(ctx, "")
	if err != nil || !dec.Allowed || dec.Remaining != 1 {
		t.Fatalf("unexpected first decision: dec=%+v err=%v", dec, err)
	}
	dec, err = l.Allow(ctx, "unknown")
	if err != nil || !dec.Allowed || dec.Remaining != 0 {
		t.Fatalf("unexpected second decision: dec=%+v err=%v", dec, err)
	}
	dec, err = l.Allow(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Allowed {
		t.Fatalf("expected blocked for empty key after shared bucket exhausted")
	}
}

func TestFixedWindowRedis_DefaultsNormalization(t *testing.T) {
	cli := checkRedisOrSkip(t)
	start := time.Now().Add(1 * time.Hour).UTC()
	storage := NewRedisStorage(cli)
	cfg := FixedWindowConfig{Window: 0, Limit: 0}
	l := NewFixedWindowLimiter(storage, cfg).(*FixedWindowLimiter)
	mc := NewMockClock(start)
	l.clock = mc

	ctx := context.Background()
	var lastDec Decision
	var err error
	for i := 0; i < 60; i++ {
		lastDec, err = l.Allow(ctx, "norm")
		if err != nil || !lastDec.Allowed {
			t.Fatalf("expected call %d allowed, got dec=%+v err=%v", i+1, lastDec, err)
		}
	}
	if lastDec.Remaining != 0 {
		t.Fatalf("expected remaining=0 after 60th allowed, got %d", lastDec.Remaining)
	}
	dec, err := l.Allow(ctx, "norm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Allowed {
		t.Fatalf("expected blocked on 61st call with default limit")
	}
	expectedReset := mc.Now().Truncate(time.Minute).Add(time.Minute)
	if !dec.ResetAt.Equal(expectedReset) {
		t.Fatalf("expected resetAt=%v, got %v", expectedReset, dec.ResetAt)
	}
}

// --- existing memory-backed tests below ---

func TestFixedWindow_AllowWithinLimitAndBlock(t *testing.T) {
	// pick a start in the near future so MemoryStorage's real time expiry doesn't invalidate entries mid-test
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(3)
	window := 10 * time.Second
	l, mc := newTestLimiter(limit, window, "", start)

	ctx := context.Background()
	// First 3 pass
	for i := int64(1); i <= limit; i++ {
		dec, err := l.Allow(ctx, "k")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !dec.Allowed {
			t.Fatalf("call %d expected allowed=true", i)
		}
		expectedRemaining := limit - i
		if dec.Remaining != expectedRemaining {
			t.Fatalf("call %d expected remaining=%d, got %d", i, expectedRemaining, dec.Remaining)
		}
		windowStart := mc.Now().Truncate(window)
		expectedReset := windowStart.Add(window)
		if !dec.ResetAt.Equal(expectedReset) {
			t.Fatalf("call %d expected resetAt=%v, got %v", i, expectedReset, dec.ResetAt)
		}
	}

	// Next one should be blocked
	dec, err := l.Allow(ctx, "k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Allowed {
		t.Fatalf("expected allowed=false after limit exceeded")
	}
	if dec.Remaining != 0 {
		t.Fatalf("expected remaining=0 when blocked, got %d", dec.Remaining)
	}
	windowStart := mc.Now().Truncate(window)
	expectedReset := windowStart.Add(window)
	if !dec.ResetAt.Equal(expectedReset) {
		t.Fatalf("expected resetAt=%v, got %v", expectedReset, dec.ResetAt)
	}
	// dec.RetryAfter is based on real time; don't assert exact value.
}

func TestFixedWindow_ResetsNextWindow(t *testing.T) {
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(3)
	window := 10 * time.Second
	l, mc := newTestLimiter(limit, window, "", start)
	ctx := context.Background()

	// Exhaust the limit in the first window
	for i := int64(0); i < limit; i++ {
		if dec, err := l.Allow(ctx, "user"); err != nil || !dec.Allowed {
			t.Fatalf("pre-window calls should be allowed: dec=%+v err=%v", dec, err)
		}
	}
	if dec, _ := l.Allow(ctx, "user"); dec.Allowed {
		t.Fatalf("expected blocked after exhausting limit")
	}

	// Advance to next window
	mc.Advance(window)
	dec, err := l.Allow(ctx, "user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dec.Allowed {
		t.Fatalf("expected allowed at the start of new window")
	}
	if dec.Remaining != limit-1 {
		t.Fatalf("expected remaining=%d, got %d", limit-1, dec.Remaining)
	}
	expectedReset := mc.Now().Truncate(window).Add(window)
	if !dec.ResetAt.Equal(expectedReset) {
		t.Fatalf("expected resetAt=%v, got %v", expectedReset, dec.ResetAt)
	}
}

func TestFixedWindow_EmptyKeyUsesUnknownBucket(t *testing.T) {
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(2)
	window := 10 * time.Second
	l, _ := newTestLimiter(limit, window, "", start)
	ctx := context.Background()

	// First call with empty key should be allowed
	dec, err := l.Allow(ctx, "")
	if err != nil || !dec.Allowed || dec.Remaining != 1 {
		t.Fatalf("unexpected first decision: dec=%+v err=%v", dec, err)
	}
	// Second call with explicit "unknown" should share the same bucket and exhaust it
	dec, err = l.Allow(ctx, "unknown")
	if err != nil || !dec.Allowed || dec.Remaining != 0 {
		t.Fatalf("unexpected second decision: dec=%+v err=%v", dec, err)
	}
	// Third call with empty key should now be blocked
	dec, err = l.Allow(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Allowed {
		t.Fatalf("expected blocked for empty key after shared bucket exhausted")
	}
}

func TestFixedWindow_KeyPrefixIsolation(t *testing.T) {
	start := time.Now().Add(1 * time.Hour).UTC()
	limit := int64(3)
	window := 10 * time.Second

	// Shared storage; separate prefixes
	storage := NewMemoryStorage()
	cfgA := FixedWindowConfig{Window: window, Limit: limit, KeyPrefix: "a:"}
	cfgB := FixedWindowConfig{Window: window, Limit: limit, KeyPrefix: "b:"}
	la := NewFixedWindowLimiter(storage, cfgA).(*FixedWindowLimiter)
	lb := NewFixedWindowLimiter(storage, cfgB).(*FixedWindowLimiter)
	mc := NewMockClock(start)
	la.clock = mc
	lb.clock = mc

	ctx := context.Background()

	// Exhaust A's window
	for i := int64(0); i < limit; i++ {
		if dec, err := la.Allow(ctx, "k"); err != nil || !dec.Allowed {
			t.Fatalf("expected la allowed, got dec=%+v err=%v", dec, err)
		}
	}
	if dec, _ := la.Allow(ctx, "k"); dec.Allowed {
		t.Fatalf("expected la blocked after limit")
	}

	// B should still be fresh and allowed
	dec, err := lb.Allow(ctx, "k")
	if err != nil || !dec.Allowed || dec.Remaining != limit-1 {
		t.Fatalf("expected lb allowed and independent, dec=%+v err=%v", dec, err)
	}
}

func TestFixedWindow_DefaultsNormalization(t *testing.T) {
	start := time.Now().Add(1 * time.Hour).UTC() // use future start
	// zero/negative config should normalize to 1m window and limit=60
	storage := NewMemoryStorage()
	cfg := FixedWindowConfig{Window: 0, Limit: 0}
	l := NewFixedWindowLimiter(storage, cfg).(*FixedWindowLimiter)
	mc := NewMockClock(start)
	l.clock = mc

	ctx := context.Background()
	// Expect exactly 60 allowed, then blocked
	var lastDec Decision
	var err error
	for i := 0; i < 60; i++ {
		lastDec, err = l.Allow(ctx, "norm")
		if err != nil || !lastDec.Allowed {
			t.Fatalf("expected call %d allowed, got dec=%+v err=%v", i+1, lastDec, err)
		}
	}
	if lastDec.Remaining != 0 {
		t.Fatalf("expected remaining=0 after 60th allowed, got %d", lastDec.Remaining)
	}
	dec, err := l.Allow(ctx, "norm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec.Allowed {
		t.Fatalf("expected blocked on 61st call with default limit")
	}
	// Ensure ResetAt is aligned to 1-minute window
	expectedReset := mc.Now().Truncate(time.Minute).Add(time.Minute)
	if !dec.ResetAt.Equal(expectedReset) {
		t.Fatalf("expected resetAt=%v, got %v", expectedReset, dec.ResetAt)
	}
}
