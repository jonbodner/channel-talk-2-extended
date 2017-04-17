package main

import (
	"errors"
	"testing"
	"time"
)

func doSomethingThatTakesAWhile(i int) (int, error) {
	// what this does isn't that important for the example
	// but it's really important in your program
	time.Sleep(2 * time.Second)
	return i * 2, nil
}

func doAnotherThing(i interface{}) (interface{}, error) {
	ii := i.(int)
	// we can wait a bit here, too
	time.Sleep(50 * time.Millisecond)
	return ii + 1, nil
}

func doAnError(i int) (int, error) {
	//this errors out
	return 0, errors.New("Nope, no good")
}

func TestThen(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	}).Then(doAnotherThing)

	// this will wait for a second to complete
	// and will return nil for both val and err
	// and timeout will be true
	val, timeout, err := f.GetUntil(1 * time.Second)
	if val != nil || timeout != true || err != nil {
		t.Errorf("Expected nil, true, nil, got %v, %v, %v", val, timeout, err)
	}

	// this will wait for val and err to have values
	val, err = f.Get()
	if val != 21 || err != nil {
		t.Errorf("Expected 21, nil, got %v, %v", val, err)
	}

	// this will error out on the first step, so we
	// don't do the long-running thing next
	g := New(func() (interface{}, error) {
		return doAnError(a)
	}).Then(func(v interface{}) (interface{}, error) {
		return doSomethingThatTakesAWhile(v.(int))
	})

	val2, timeout2, err2 := g.GetUntil(1 * time.Second)
	if val2 != 0 || timeout2 != false || err2 == nil || err2.Error() != "Nope, no good" {
		t.Errorf("Expected 0, false, Nope, no good, got %v, %v, %v", val2, timeout2, err2)
	}
}
