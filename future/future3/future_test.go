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

func TestGetUntil(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	})

	// this will wait about a second to complete
	// and will return nil for both val and err
	// and timeout will be true
	start := time.Now()
	val, timeout, err := f.GetUntil(time.Second)
	end := time.Now()
	if end.Sub(start) < time.Second {
		t.Errorf("That should have waited a second to finish, took %v", end.Sub(start))
	}
	if val != nil || err != nil || timeout != true {
		t.Errorf("Expected nil, true, and nil, got %v, %v, and %v", val, timeout, err)
	}
	fmt.Println(val, err)

	// this will wait for val and err to have values
	start = time.Now()
	val, err = f.Get()
	end = time.Now()
	if end.Sub(start) > 1100*time.Millisecond {
		t.Errorf("That should have taken about a second to finish, took %v", end.Sub(start))
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}
	fmt.Println(val, err)

	//this GetUntil will return immediately
	start = time.Now()
	val, timeout, err = f.GetUntil(time.Second)
	end = time.Now()
	if end.Sub(start) > 100*time.Millisecond {
		t.Errorf("That should have taken no time at all to finish, took %v", end.Sub(start))
	}
	if val != 20 || err != nil || timeout {
		t.Errorf("Expected 20, false, and nil, got %v, %v and %v", val, timeout, err)
	}
	fmt.Println(val, err)
}