package zcache

import (
	"sync"
	"testing"
	"time"
)

var (
	once        sync.Once
	stringGroup *Group
	stringc     = make(chan string)
)

const (
	stringGroupName = "string-group"
	fromChan        = "from-chan"
	cacheSize       = 1 << 20
)
func testSetup() {
	stringGroup = NewGroup(stringGroupName, cacheSize, GetterFunc(func(key string) ([]byte, error) {
		if key == fromChan {
			key = <-stringc
		}
		val := []byte("ECHO:" + key)
		return val, nil
	}))
}

func TestGetDupSuppressString(t *testing.T) {
	once.Do(testSetup)
	resc := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			var (
				val ByteView
				err error
			)
			val, err = stringGroup.Get(fromChan)
			if err != nil {
				resc <- "ERROR:" + err.Error()
				return
			}
			resc <- string(val.ByteSlice())
		}()
	}

	time.Sleep(250 * time.Millisecond)

	stringc <- "foo"

	for i := 0; i < 2; i++ {
		select {
		case v := <-resc:
			if v != "ECHO:foo" {
				t.Errorf("got %q; want %q", v, "ECHO:foo")
			}
		case <-time.After(5 * time.Second):
			t.Errorf("timeout waiting on getter #%d of 2", i+1)
		}
	}
}
