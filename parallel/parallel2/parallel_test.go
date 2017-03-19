package main

import (
	"fmt"
	"testing"
	"time"
)

func TestFanoutTimed(t *testing.T) {
	type in struct {
		a int
		b int
	}
	type out struct {
		source string
		result int
	}
	evaluators := []Evaluator{
		func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a + i.b
			return out{"adder", r}, nil
		},
		func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a * i.b
			return out{"timer", r}, nil
		},
		func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a - i.b
			return out{"subber", r}, nil
		},
		func(inV interface{}) (interface{}, error) {
			time.Sleep(20*time.Millisecond)
			i := inV.(in)
			r := i.a / i.b
			return out{"divider", r}, nil
		},
	}
	results, _ := FanOutTimed(in{2, 3}, evaluators, 10*time.Millisecond)
	fmt.Println(results)
}
