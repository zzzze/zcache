package zcache

import (
	"log"
	"sync"
	"zcache/singleflight"
	pb "zcache/zcachepb"
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
  peers PeerPicker
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
  if val, ok := g.mainCache.get(key); ok {
    log.Println("[zcache] hit")
    return val, nil
  }
  return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
  val, err := g.loader.Do(key, func() (interface{}, error) {
    if g.peers != nil {
      if peer, ok := g.peers.PickPeer(key); ok {
        var res pb.Response
        req := &pb.Request{
          Group: g.name,
          Key: key,
        }
        if err := peer.Get(req, &res); err != nil {
          return ByteView{}, err
        }
        return ByteView{res.Value}, nil
      }
    }
    return g.getLocally(key)
  })
  return val.(ByteView), err
}

func (g *Group) getLocally(key string) (ByteView, error) {
  data, err := g.getter.Get(key)
  v := ByteView{data}
  if err != nil {
    return v, err
  }
  g.mainCache.add(key, v)
  return v, nil
}

// RegisterPeers ...
func (g *Group) RegisterPeers(peers PeerPicker) {
  if g.peers != nil {
    panic("RegisterPeers called more then once")
  }
  g.peers = peers
}
