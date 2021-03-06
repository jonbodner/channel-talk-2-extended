package main

import (
	"sync"
	"testing"
	"fmt"
)

func TestPool(t *testing.T) {
	adderPool := NewPool(adderFactory, 3)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				//borrow an adder
				adder := adderPool.Borrow().(func(int, int) int)
				fmt.Println("adding on goroutine", i)
				fmt.Println("result == ", adder(i, i + 1))
				//return the item
				adderPool.Return(adder)
			}
		}(i)
	}
	wg.Wait()
}
