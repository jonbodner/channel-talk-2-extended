package main

import (
	"fmt"
	"testing"
	"time"
)

func doSomethingThatTakesAWhile(i int) (int, error) {
	// what this does isn't that important for the example
	// but it's really important in your program
	time.Sleep(2 * time.Second)
	return i * 2, nil
}

func TestGet(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	})

	// this will take about 2 seconds to complete
	start := time.Now()
	val, err := f.Get()
	end := time.Now()
	if end.Sub(start) < time.Second {
		t.Errorf("That should have taken 2 seconds to finish, took %v", end.Sub(start))
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}
	fmt.Println(val, err)

	// this will complete immediately
	start = time.Now()
	val, err = f.Get()
	end = time.Now()
	if end.Sub(start) > time.Second {
		t.Errorf("That should have taken no time to finish, took %v", end.Sub(start))
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}
	fmt.Println(val, err)
}
