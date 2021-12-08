package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
)

// Hash ...
type Hash func(data []byte) uint32

// Map ...
type Map struct {
	hash     Hash
	replicas int            // 虚拟节点个数
	keys     []int          // 哈希环
	hashMap  map[int]string // 节点哈希于节点的对应关系
}

// New ...
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加节点到环上
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas+1; i++ {
			h := int(m.hash([]byte(fmt.Sprintf("%d%s", i, key))))
			m.keys = append(m.keys, h)
			m.hashMap[h] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 获取命中的节点
func (m *Map) Get(key string) string {
  if len(m.keys) == 0 {
    return ""
  }
  h := int(m.hash([]byte(key)))
  idx := sort.Search(len(m.keys), func(i int) bool {
    return m.keys[i] >= h
  })
  return m.hashMap[m.keys[idx%len(m.keys)]]
}
