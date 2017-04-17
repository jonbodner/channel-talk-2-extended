package main

import (
	"fmt"
	"testing"
	"time"
)

type in struct {
	a int
	b int
}

type out struct {
	source string
	result int
}

func adder(i in) (out, error) {
	r := i.a + i.b
	return out{"adder", r}, nil
}

func timer(i in) (out, error) {
	r := i.a * i.b
	return out{"timer", r}, nil
}

func subber(i in) (out, error) {
	r := i.a - i.b
	return out{"subber", r}, nil
}

func divider(i in) (out, error) {
	time.Sleep(20*time.Millisecond)
	r := i.a / i.b
	return out{"divider", r}, nil
}

func TestFanoutTimed(t *testing.T) {
	evaluators := []Evaluator{
		func(inV interface{}) (interface{}, error) {
			return adder(inV.(in))
		},
		func(inV interface{}) (interface{}, error) {
			return timer(inV.(in))
		},
		func(inV interface{}) (interface{}, error) {
			return subber(inV.(in))
		},
		func(inV interface{}) (interface{}, error) {
			return divider(inV.(in))
		},
	}
	results, errors := FanOutTimed(in{2, 3}, evaluators, 10 * time.Millisecond)
	fmt.Println("results:", results)
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
	fmt.Println("errors:", errors)
	if len(errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errors))
	}
}
