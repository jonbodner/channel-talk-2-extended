package main

import "time"

type Evaluator func(interface{}) (interface{}, error)

func FanOutTimed(data interface{}, evaluators []Evaluator, timeout time.Duration) ([]interface{}, []error) {
	results, errors := launch(data, evaluators)
	out, errs := gather(results, errors, len(evaluators), timeout)
	return out, errs
}

func launch(data interface{}, evaluators []Evaluator) (chan interface{}, chan error) {
	results := make(chan interface{}, len(evaluators))
	errors := make(chan error, len(evaluators))
	for _, v := range evaluators {
		go func(e Evaluator) {
			result, err := e(data)
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}(v)
	}
	return results, errors
}

func gather(results chan interface{}, errors chan error, count int, timeout time.Duration) ([]interface{}, []error) {
	out := make([]interface{}, 0, count)
	errs := make([]error, 0, count)
	timer := time.After(timeout)
	loop:
	for i := 0; i < count; i++ {
		select {
		case r := <-results:
			out = append(out, r)
		case e := <-errors:
			errs = append(errs, e)
		case <- timer:
			break loop
		}
	}
	return out, errs
}
