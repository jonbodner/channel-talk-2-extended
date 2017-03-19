package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p := NewPool(func() interface{} {
		return 5
	}, 3)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				v := p.Borrow().(int)
				fmt.Println(i, i+v)
				time.Sleep(time.Duration((i+1)*100) * time.Millisecond)
				p.Return(v)
			}
		}(i)
	}
	wg.Wait()
}
