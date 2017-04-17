package main

import (
	"time"
)

type Future interface {
	Get() (interface{}, error)

	GetUntil(d time.Duration) (interface{}, bool, error)

	Then(Step) Future

	Cancel()

	IsCancelled() bool
}

type Process func() (interface{}, error)

type Step func(interface{}) (interface{}, error)

type futureImpl struct {
	done   chan struct{}
	cancel chan struct{}
	guard  chan struct{}
	val    interface{}
	err    error
}

func (f *futureImpl) Get() (interface{}, error) {
	select {
	case <-f.done:
		return f.val, f.err
	case <-f.cancel:
	//on cancel, just fall out
	}
	return nil, nil
}

func (f *futureImpl) GetUntil(d time.Duration) (interface{}, bool, error) {
	select {
	case <-f.done:
		val, err := f.Get()
		return val, false, err
	case <-time.After(d):
		return nil, true, nil
	case <-f.cancel:
	//on cancel, just fall out
	}
	return nil, false, nil
}

func (f *futureImpl) Then(next Step) Future {
	nextFuture := newInner(f.cancel, f.guard, func() (interface{}, error) {
		result, err := f.Get()
		if f.IsCancelled() || err != nil {
			return result, err
		}
		return next(result)
	})
	return nextFuture
}

func (f *futureImpl) Cancel() {
	select {
	case <-f.done:
	//already finished
	case <-f.cancel:
	//already cancelled
	default:
	//making request to cancel
	// this will serialize requests from multiple goroutines
			select {
			case <-f.done:
			case <-f.cancel:
			case f.guard <- struct{}{}:
			}
	}
}

func (f *futureImpl) IsCancelled() bool {
	select {
	case <-f.cancel:
		return true
	default:
		return false
	}
}

func New(inFunc Process) Future {
	cancel := make(chan struct{})
	guard := make(chan struct{})
	go func() {
		<-guard
		close(cancel)
	}()
	return newInner(cancel, guard, inFunc)
}

func newInner(cancel chan struct{}, guard chan struct{}, inFunc Process) Future {
	f := futureImpl{
		done:   make(chan struct{}),
		cancel: cancel,
		guard: guard,
	}
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return &f
}
