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

func NewCache(interval time.Duration) *Cache {
	cache := Cache{}
	cache.interval = interval
	cache.mu = sync.Mutex{}
	cache.entry = map[string]cacheEntry{}
	go cache.reapLoop()
	return &cache
}

func (cache *Cache) Add(key string, val []byte) {
	newEntry := cacheEntry{}
	newEntry.createdAt = time.Now()
	newEntry.val = val
	cache.mu.Lock()
	cache.entry[key] = newEntry
	cache.mu.Unlock()
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	val, ok := cache.entry[key]
	if ok {
		return val.val, true
	} else {
		return nil, false
	}
}

func (cache *Cache) reapLoop() {
	for {
		cache.mu.Lock()
		for key, val := range cache.entry {
			current := time.Now()
			if current.Sub(val.createdAt) > cache.interval {
				delete(cache.entry, key)	
			}
		}
		cache.mu.Unlock()
		time.Sleep(cache.interval)
	}
}
