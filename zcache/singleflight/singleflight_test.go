package singleflight

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleFlight(t *testing.T) {
	data := make(chan string)
	count := 10
	var (
		g  Group
		wg sync.WaitGroup
	)
	calls := int64(0)
	fn := func() (interface{}, error) {
		atomic.AddInt64(&calls, 1)
		return <-data, nil
	}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			d, err := g.Do("myKey", fn)
			if err != nil {
				t.Errorf("got error: %q", err.Error())
			}
			if d.(string) != "data" {
				t.Errorf("got %q, want %q", d, "data")
			}
			wg.Done()
		}()
	}
	time.Sleep(time.Millisecond * 100)
	data <- "data"
	wg.Wait()
	if got := atomic.LoadInt64(&calls); got != 1 {
		t.Errorf("number of calls = %d times, want 1", got)
	}
}
