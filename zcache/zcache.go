package zcache

import (
	"sync"
	"zcache/singleflight"
)

var mu sync.RWMutex
var groups = make(map[string]*Group)

// Getter ...
type Getter interface {
  Get(key string) ([]byte, error)
}

// GetterFunc ...
type GetterFunc func(key string) ([]byte, error)

// Get ...
func (f GetterFunc) Get(key string) ([]byte, error) {
  return f(key)
}

// GetGroup ...
func GetGroup(name string) *Group {
  mu.RLock()
  defer mu.RUnlock()
  return groups[name]
}

// Group ...
type Group struct {
  name string
  mainCache cache
  cacheBytes int64
  getter Getter
  loader singleflight.Group
}

// NewGroup ...
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
  mu.Lock()
  defer mu.Unlock()
  g := &Group{
    name: name,
    getter: getter,
    cacheBytes: cacheBytes,
    mainCache: cache{maxBytes: cacheBytes},
  }
  groups[name] = g
  return g
}

// Get ...
func (g *Group) Get(key string) (ByteView, error) {
  val, err := g.loader.Do(key, func() (interface{}, error) {
    val, ok := g.mainCache.get(key)
    if !ok {
      data, err := g.getter.Get(key)
      v := ByteView{data}
      if err != nil {
        return v, err
      }
      g.mainCache.add(key, v)
      return v, nil
    }
    return val, nil
  })
  return val.(ByteView), err
}
