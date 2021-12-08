package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		i, err := strconv.Atoi(string(data))
		if err != nil {
			t.Error(err.Error())
		}
		return uint32(i)
	})
	// 2, 4, 6, 12, 14, 16, 22, 24, 26, 32, 34, 36
	hash.Add("6", "4", "2")
	testCases := map[string]string {
		"9": "2",
		"16": "6",
		"33": "4",
		"38": "2",
	}
	for k, v := range testCases {
		val := hash.Get(k)
		if val != v {
			t.Errorf("asking for %q, got %q, want %q", k, val, v)
		}
	}
	hash.Add("8")
	testCases2 := map[string]string {
		"9": "2",
		"16": "6",
		"33": "4",
		"38": "8",
	}
	for k, v := range testCases2 {
		val := hash.Get(k)
		if val != v {
			t.Errorf("asking for %q, got %q, want %q", k, val, v)
		}
	}
}
