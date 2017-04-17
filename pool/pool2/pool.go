package main

import (
	"fmt"
	"time"
)

type Factory func() interface{}

type Processor func(interface{})

type Pool interface {
	Run(Processor)
	RunWithTimeout(Processor, time.Duration) error
}

type poolInner struct {
	items chan interface{}
}

func NewPool(f Factory, count int) Pool {
	pI := &poolInner{items: make(chan interface{}, count)}
	for i := 0; i < count; i++ {
		pI.items <- f()
	}
	return pI
}

func (pi *poolInner) Run(p Processor) {
	item := <-pi.items
	defer func() {
		pi.items <- item
	}()
	p(item)
}

func (pi *poolInner) RunWithTimeout(p Processor, d time.Duration) error {
	select {
	case item := <-pi.items:
		defer func() {
			pi.items <- item
		}()
		p(item)
	case <-time.After(d):
		return fmt.Errorf("Pool timed out after %s", d.String())
	}
	return nil
}
