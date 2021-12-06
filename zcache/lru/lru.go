package lru

import "container/list"

// Cache ...
type Cache struct {
	MaxEntries int
	ll         *list.List
	cache      map[interface{}]*list.Element
	OnEvicted  func(key, value interface{})
}

type entry struct {
	key   interface{}
	value interface{}
}

// New ...
func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

// Add ...
func (c *Cache) Add(key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if item, hit := c.cache[key]; hit {
		c.ll.MoveToFront(item)
		item.Value.(*entry).value = value
		return
	}
	item := c.ll.PushFront(&entry{key, value})
	c.cache[key] = item
  if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
    c.RemoveOldest()
  }
}

// Get ...
func (c *Cache) Get(key interface{}) (value interface{}, hit bool) {
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
func (c *Cache) Remove(key interface{}) {
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
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}
