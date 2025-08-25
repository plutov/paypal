package limiter

import (
	"context"
	"hash/fnv"
	"strconv"
	"sync"
	"time"
)

type MemoryStorage struct {
	shards          []*memoryShard
	shardMask       uint32
	cleanupInterval time.Duration
}

type memoryShard struct {
	mu        sync.RWMutex
	data      map[string]memoryValue
	swHits    map[string][]int64
	swBlocked map[string]time.Time
}

type memoryValue struct {
	value  string
	expiry time.Time
}

func NewMemoryStorage() FixedWindowBackend { return NewMemoryStorageWithShards(16) }

func NewMemoryStorageWithShards(shards int) FixedWindowBackend {
	if shards <= 0 {
		shards = 16
	}

	s := 1
	for s < shards {
		s <<= 1
	}
	ms := &MemoryStorage{shards: make([]*memoryShard, s), shardMask: uint32(s - 1), cleanupInterval: time.Minute}
	for i := 0; i < s; i++ {
		ms.shards[i] = &memoryShard{data: make(map[string]memoryValue), swHits: make(map[string][]int64), swBlocked: make(map[string]time.Time)}
	}
	go ms.cleanup()
	return ms
}

func (m *MemoryStorage) IncrWindow(ctx context.Context, baseKey string, window time.Duration, now time.Time) (int64, time.Time, error) {
	key := baseKey + ":" + strconv.FormatInt(now.Truncate(window).UnixMilli(), 10)
	sh := m.pick(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	windowStart := now.Truncate(window)
	resetAt := windowStart.Add(window)
	mv := sh.data[key]
	if !mv.expiry.IsZero() && time.Now().After(mv.expiry) {
		mv.value = "0"
	}
	cnt, _ := strconv.ParseInt(mv.value, 10, 64)
	cnt++
	mv.value = strconv.FormatInt(cnt, 10)
	mv.expiry = resetAt.Add(time.Second)
	sh.data[key] = mv
	return cnt, resetAt, nil
}

func (m *MemoryStorage) pick(key string) *memoryShard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return m.shards[h.Sum32()&m.shardMask]
}

func (m *MemoryStorage) cleanup() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		for _, sh := range m.shards {
			sh.mu.Lock()
			for k, v := range sh.data {
				if !v.expiry.IsZero() && now.After(v.expiry) {
					delete(sh.data, k)
				}
			}

			for k, until := range sh.swBlocked {
				if now.After(until) {
					delete(sh.swBlocked, k)
				}
			}
			sh.mu.Unlock()
		}
	}
}
