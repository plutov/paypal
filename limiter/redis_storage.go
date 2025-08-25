package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const luaFixedWindow = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local windowStart = math.floor(now / window) * window
local windowKey = key .. ":" .. windowStart

local current = redis.call('INCR', windowKey)
if current == 1 then
  redis.call('PEXPIRE', windowKey, window + 1000)
end
local resetTime = windowStart + window
return {current, resetTime}
`

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) FixedWindowBackend {
	return &RedisStorage{client: client}
}

func (r *RedisStorage) IncrWindow(ctx context.Context, baseKey string, window time.Duration, now time.Time) (int64, time.Time, error) {
	res, err := r.client.Eval(ctx, luaFixedWindow, []string{baseKey}, now.UnixMilli(), window.Milliseconds()).Result()
	if err != nil {
		return 0, time.Time{}, err
	}
	arr, ok := res.([]interface{})
	if !ok || len(arr) != 2 {
		return 0, time.Time{}, ErrUnexpectedScriptResult
	}
	count := arr[0].(int64)
	resetMs := arr[1].(int64)
	return count, time.UnixMilli(resetMs), nil
}
