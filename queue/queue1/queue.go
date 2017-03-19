package main

func MakeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	//do magic stuff 
	return in, out
}


