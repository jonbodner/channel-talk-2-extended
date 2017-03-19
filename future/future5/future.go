package main

import (
	"time"
)

type Interface interface {
	Get() (interface{}, error)

	GetUntil(d time.Duration) (interface{}, bool, error)

	Then(Step) Interface

	Cancel()

	IsCancelled() bool
}

type Process func() (interface{}, error)

type Step func(interface{}) (interface{}, error)

type futureImpl struct {
	done   chan struct{}
	cancel chan struct{}
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

func (f *futureImpl) Then(next Step) Interface {
	nextFuture := newInner(f.cancel, func() (interface{}, error) {
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
		close(f.cancel) //should only be called once, since the closed cancel channel will always return
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

func New(inFunc Process) Interface {
	return newInner(make(chan struct{}), inFunc)
}

func newInner(cancelChan chan struct{}, inFunc Process) Interface {
	f := futureImpl{
		done:   make(chan struct{}),
		cancel: cancelChan,
	}
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return &f
}
