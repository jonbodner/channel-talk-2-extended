package main

func MakeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
		curVal := func() interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return inQueue[0]
		}
	loop:
		for {
			select {
			case v, ok := <-in:
				if !ok {
					break loop
				} else {
					inQueue = append(inQueue, v)
				}
			case out <- curVal():
				if len(inQueue) > 0 {
					inQueue = inQueue[1:]
				}
			}
		}
		close(out)
	}()
	return in, out
}
