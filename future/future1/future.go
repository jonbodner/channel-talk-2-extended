package main

type Interface interface {
	Get() (interface{}, error)
}

type Process func() (interface{}, error)

type futureImpl struct {
	//this will have stuff eventually
}

func (f *futureImpl) Get() (interface{}, error) {
	//this will have stuff eventually
	return nil, nil
}

func New(inFunc Process) Interface {
	//do magic stuff
	return &futureImpl{}
}

