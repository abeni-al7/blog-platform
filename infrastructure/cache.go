package infrastructure

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
	hasExpiry bool
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]cacheItem
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]cacheItem)}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := cacheItem{value: value}
	if ttl > 0 {
		item.hasExpiry = true
		item.expiresAt = time.Now().Add(ttl)
	}
	c.data[key] = item
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if item.hasExpiry && time.Now().After(item.expiresAt) {
		c.Delete(key)
		return nil, false
	}
	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *Cache) Clear() {
	c.mu.Lock()
	c.data = make(map[string]cacheItem)
	c.mu.Unlock()
}
