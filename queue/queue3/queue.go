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
			case out <- inQueue[0]:
				inQueue = inQueue[1:]
			}
		}
		close(out)
	}()
	return &queueInner{in, out}
}
