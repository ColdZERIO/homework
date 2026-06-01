package storage

import "sync"

type MemoryCache struct {
	mu    sync.Mutex
	cache map[string]any
}

func UserMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string]any),
	}
}

func (mc *MemoryCache) Get(key string) (any, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	value, ok := mc.cache[key]
	return value, ok
}

func (mc *MemoryCache) Set(key string, value any) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache[key] = value
}

func (mc *MemoryCache) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.cache, key)
}

func (mc *MemoryCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cache = make(map[string]any)
}
