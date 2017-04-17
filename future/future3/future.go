package main

import "time"

type Process func() (interface{}, error)

func New(inFunc Process) Future {
	f := &futureImpl{
		done: make(chan struct{}),
	}
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return f
}

type Future interface {
	Get() (interface{}, error)

	GetUntil(d time.Duration) (interface{}, bool, error)
}

type futureImpl struct {
	done chan struct{}
	val  interface{}
	err  error
}

func (f *futureImpl) Get() (interface{}, error) {
	<-f.done
	return f.val, f.err
}

func (f *futureImpl) GetUntil(d time.Duration) (interface{}, bool, error) {
	select {
	case <-f.done:
		val, err := f.Get()
		return val, false, err
	case <-time.After(d):
		return nil, true, nil
	}
	// This should never be executed
	return nil, false, nil
}

