package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	value          string
	expirationTime int64
}

func (c cacheItem) isExpired() bool {
	return c.expirationTime != 0 && c.expirationTime < time.Now().UnixNano()
}

type Cache struct {
	lock  *sync.RWMutex
	items map[string]cacheItem
}

func NewCache() Cache {
	return Cache{
		lock:  &sync.RWMutex{},
		items: map[string]cacheItem{},
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.lock.RLock()
	defer c.lock.Unlock()

	item, ok := c.items[key]
	if !ok || item.isExpired() {
		return "", false
	}

	return item.value, true
}

func (c *Cache) Put(key, value string) {
	c.PutTill(key, value, time.Time{})
}

func (c *Cache) Keys() []string {
	c.lock.RLock()
	defer c.lock.Unlock()

	result := make([]string, 0)
	for key, item := range c.items {
		if !item.isExpired() {
			result = append(result, key)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = cacheItem{
		value:          value,
		expirationTime: deadline.UnixNano(),
	}
}
