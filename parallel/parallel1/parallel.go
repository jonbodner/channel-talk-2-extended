package main

type Evaluator func(interface{}) (interface{}, error)

func FanOut(evaluators []Evaluator, data interface{}) ([]interface{}, []error) {
	results, errors := launch(evaluators, data)
	out, errs := gather(results, errors, len(evaluators))
	return out, errs
}

func launch(evaluators []Evaluator, data interface{}) (chan interface{}, chan error) {
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

func gather(results chan interface{}, errors chan error, count int) ([]interface{}, []error) {
	out := make([]interface{}, 0, count)
	errs := make([]error, 0, count)
	for i := 0; i < count; i++ {
		select {
		case r := <-results:
			out = append(out, r)
		case e := <-errors:
			errs = append(errs, e)
		}
	}
	return out, errs
}
