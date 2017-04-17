package main

import (
	"sync"
	"testing"
	"time"
	"sync/atomic"
	"fmt"
)

func itemFactory() interface{} {
	return struct{}{}
}

func TestPool(t *testing.T) {
	p := NewPool(itemFactory, 3)

	var wg sync.WaitGroup
	var count int32
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				//try to borrow something
				v := p.Borrow()
				//increment the count of borrowed items
				atomic.AddInt32(&count, 1)
				//show that it is never greater than 3
				fmt.Println(count)
				if count > 3 {
					t.Errorf("Should never have more than 3 goroutines active at once, but have %d", count)
				}
				//sleep for a bit
				time.Sleep(time.Duration((i + 1) * 100) * time.Millisecond)
				//decrement the count of borrowed items
				atomic.AddInt32(&count, -1)
				//return the item
				p.Return(v)
			}
		}(i)
	}
	wg.Wait()
}
