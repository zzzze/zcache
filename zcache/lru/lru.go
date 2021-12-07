package lru

import "container/list"

// Cache ...
type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[interface{}]*list.Element
	OnEvicted func(key, value interface{})
}

// Value ...
type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

// New ...
func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[interface{}]*list.Element),
	}
}

// Add ...
func (c *Cache) Add(key string, value Value) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if item, hit := c.cache[key]; hit {
		c.ll.MoveToFront(item)
		kv := item.Value.(*entry)
		c.nBytes += int64(value.Len() - kv.value.Len())
		kv.value = value
		return
	}
	item := c.ll.PushFront(&entry{key, value})
	c.cache[key] = item
	c.nBytes += int64(len(key) + value.Len())
	if c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Get ...
func (c *Cache) Get(key string) (value Value, hit bool) {
	if c.cache == nil {
		return nil, false
	}
	item, hit := c.cache[key]
	if !hit {
		return nil, false
	}
	c.ll.MoveToFront(item)
	return item.Value.(*entry).value, true
}

// Remove ...
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	item, hit := c.cache[key]
	if !hit {
		return
	}
	c.removeElement(item)
}

// RemoveOldest ...
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	item := c.ll.Back()
	c.removeElement(item)
}

func (c *Cache) removeElement(item *list.Element) {
	c.ll.Remove(item)
	kv := item.Value.(*entry)
	c.nBytes -= int64(len(kv.key) + kv.value.Len())
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}
