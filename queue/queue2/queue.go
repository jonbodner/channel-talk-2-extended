package main

func MakeInfiniteQueue() Queue {
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
			// do something here that makes
			// a value go on to the out channel
			}
		}
		close(out)
	}()
	return &queueInner{in: in, out: out}
}
