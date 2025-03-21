package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mutex sync.RWMutex
	cache map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	cache := new(Cache)
	cache.cache = make(map[string]cacheEntry, 0)
	go cache.reapLoop(interval)
	return cache
}

func (cache *Cache) Add(key string, val []byte) error {
	cache.mutex.Lock()
	cache.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.mutex.Unlock()
	return nil
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	var val []byte
	cache.mutex.RLock()
	entry, ok := cache.cache[key]
	if ok {
		val = entry.val
	}
	cache.mutex.RUnlock()
	return val, ok
}

func (cache *Cache) reapLoop(interval time.Duration) {
	for range time.Tick(interval) {
		//fmt.Println("Tick...")
		cache.mutex.Lock()
		for key, entry := range cache.cache {
			duration := time.Now().Sub(entry.createdAt)
			if duration > interval {
				fmt.Printf("Deleting cache entry '%s'...\n", key)
				delete(cache.cache, key)
			}
		}
		cache.mutex.Unlock()
	}
}
