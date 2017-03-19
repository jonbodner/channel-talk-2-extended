package main

func MakeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
	loop:
		for {
			select {
			case v, ok := <-in:
				if !ok {
					break loop
				} else {
					inQueue = append(inQueue, v)
				}
			case out <- inQueue[0]:
				inQueue = inQueue[1:]
			}
		}
		close(out)
	}()
	return in, out
}
