package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Evaluator interface {
	Evaluate(interface{}) (interface{}, error)
	Name() string
}

type EvaluatorFunc func(interface{}) (interface{}, error)

func (ef EvaluatorFunc) Evaluate(in interface{}) (interface{}, error) {
	return ef(in)
}

func (ef EvaluatorFunc) Name() string {
	return runtime.FuncForPC(reflect.ValueOf(ef).Pointer()).Name()
}

type evaluatorInner struct {
	EvaluatorFunc
	name string
}

func (ei evaluatorInner) Name() string {
	return ei.name
}

func Name(ef EvaluatorFunc, name string) Evaluator {
	return evaluatorInner{
		EvaluatorFunc: ef,
		name:          name,
	}
}

func FanOutTimedNamed(data interface{}, evaluators []Evaluator, timeout time.Duration) ([]interface{}, []error) {

	type wrapper struct {
		name   string
		result interface{}
	}

	type wrapperErr struct {
		name string
		err  error
	}

	gather := make(chan wrapper, len(evaluators))
	errors := make(chan wrapperErr, len(evaluators))

	names := map[string]bool{}

	for _, v := range evaluators {
		name := v.Name()
		names[name] = true
		go func(e Evaluator) {
			result, err := e.Evaluate(data)
			if err != nil {
				errors <- wrapperErr{name: name, err: err}
			} else {
				gather <- wrapper{name: name, result: result}
			}
		}(v)
	}
	out := make([]interface{}, 0, len(evaluators))
	errs := make([]error, 0, len(evaluators))
	timer := time.After(timeout)
loop:
	for range evaluators {
		select {
		case r := <-gather:
			out = append(out, r.result)
			delete(names, r.name)
		case e := <-errors:
			errs = append(errs, e.err)
			delete(names, e.name)
		case <-timer:
			break loop
		}
	}
	for name := range names {
		errs = append(errs, fmt.Errorf("%s timed out after %v on %v", name, timeout, data))
	}
	return out, errs
}
