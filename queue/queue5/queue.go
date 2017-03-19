package main

func MakeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
		outCh := func() chan interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return out
		}
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
			case outCh() <- curVal():
				inQueue = inQueue[1:]
			}
		}
		close(out)
	}()
	return in, out
}
