package zcache

import (
	"sync"
	"zcache/lru"
)

type cache struct {
	maxBytes int64
	mu       sync.Mutex
	lru      *lru.Cache
}

func (c *cache) add(key string, val ByteView) {
  c.mu.Lock()
  defer c.mu.Unlock()
  if c.lru == nil {
    c.lru = lru.New(c.maxBytes)
  }
  c.lru.Add(key, val)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
  c.mu.Lock()
  defer c.mu.Unlock()
  if c.lru == nil {
    return
  }
  if val, ok := c.lru.Get(key); ok {
    return val.(ByteView), ok
  }
  return
}
