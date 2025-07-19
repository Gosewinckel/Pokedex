package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	entry map[string]cacheEntry
	mu sync.Mutex
	interval time.Duration
}

func newCache(interval time.Duration) Cache {
	val cache = Cache{}
	cache.interval = interval
	return cache
}

func (cache Cache) Add(key string, val []byte) {
	cache.entry[key] := val
}
