package main

import (
	"errors"
	"fmt"
	"time"
)

type Factory func() interface{}

type poolInner struct {
	items chan interface{}
}

type Pool interface {
	Borrow() interface{}
	Return(interface{})
	BorrowWithTimeout(time.Duration) (interface{}, error)
}

func NewPool(f Factory, count int) Pool {
	pI := &poolInner{items: make(chan interface{}, count)}
	for i := 0; i < count; i++ {
		pI.items <- f()
	}
	return pI
}

func (pi *poolInner) Borrow() interface{} {
	item := <-pi.items
	return item
}

func (pi *poolInner) Return(in interface{}) {
	pi.items <- in
}

func (pi *poolInner) BorrowWithTimeout(d time.Duration) (interface{}, error) {
	select {
	case item := <-pi.items:
		return item, nil
	case <-time.After(d):
		return nil, fmt.Errorf("Pool timed out after %s", d.String())
	}
	return nil, errors.New("Should never get here!")
}
