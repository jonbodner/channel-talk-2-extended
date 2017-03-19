package main

import (
	"fmt"
	"testing"
	"time"
)

func TestFanoutTimedNamed(t *testing.T) {
	type in struct {
		a int
		b int
	}
	type out struct {
		source string
		result int
	}
	evaluators := []Evaluator{
		EvaluatorFunc(func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a + i.b
			return out{"adder", r}, nil
		}),
		Name(func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a * i.b
			time.Sleep(50 * time.Millisecond)
			return out{"timer", r}, nil
		}, "timer"),
		EvaluatorFunc(func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a - i.b
			return out{"subber", r}, nil
		}),
		EvaluatorFunc(func(inV interface{}) (interface{}, error) {
			i := inV.(in)
			r := i.a / i.b
			time.Sleep(30 * time.Millisecond)
			return out{"divider", r}, nil
		}),
	}
	results, errors := FanOutTimedNamed(in{2, 3}, evaluators, 10*time.Millisecond)
	fmt.Println(results)
	fmt.Println(errors)
}
