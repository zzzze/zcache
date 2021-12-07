package lru

import (
	"encoding/binary"
	"testing"
)

type value int

func (v value) Len() int {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(v))
	return len(buf[:n])
}

var getTests = []struct{
	name string
	keyToAdd string
	keyToGet string
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
}


func TestGet(t *testing.T) {
	for _, tt := range getTests {
		lru := New(0)
		lru.Add(tt.keyToAdd, value(1234))
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != value(1234) {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestRemove(t *testing.T) {
	lru := New(0)
	lru.Add("myKey", value(1234))
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatalf("TestRemove returned no matched")
	} else if val != value(1234) {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}
	lru.Remove("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatalf("TestRemove returned a removed entry")
	}
}
