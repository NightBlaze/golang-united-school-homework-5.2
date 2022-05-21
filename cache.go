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
	sync.RWMutex
	items map[string]cacheItem
}

func NewCache() Cache {
	return Cache{
		items: map[string]cacheItem{},
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.RLock()
	defer c.Unlock()

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
	c.RLock()
	defer c.Unlock()

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
	c.Lock()
	defer c.Unlock()

	c.items[key] = cacheItem{
		value:          value,
		expirationTime: deadline.UnixNano(),
	}
}
