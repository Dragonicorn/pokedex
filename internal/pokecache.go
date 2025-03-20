package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	map[string]cacheEntry
}

func NewCache(interval time.Duration) Cache {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	tick := make(chan bool)
	go func() {
		time.Sleep(interval)
		tick <- true
	}()
	for {
		select {
		case <-tick:
			fmt.Println("Tick...")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
		}
	}
	return make(Cache)
}

func (cache *Cache) Add(key string, val []byte) error {
	cache[key] = cacheEntry{
		createdAt: time.Now()
		val: val
	}
	return nil
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	var val []byte
	entry, ok := cache[key]
	if ok {
		val := entry.val
	}
	return val, ok
}

func (cache *Cache) reapLoop() {
	fmt.Println("reapLoop...")
}