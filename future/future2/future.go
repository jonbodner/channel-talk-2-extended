package main

type Interface interface {
	Get() (interface{}, error)
}

type Process func() (interface{}, error)

type futureImpl struct {
	done chan struct{}
	val interface{}
	err error
}

func (f *futureImpl) Get() (interface{}, error) {
	<- f.done
	return f.val, f.err
}

func New(inFunc Process) Interface {
	f := &futureImpl {
		done: make(chan struct{}),
	}
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return f
}

