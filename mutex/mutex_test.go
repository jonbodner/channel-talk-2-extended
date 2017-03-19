package mutex_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	//"github.com/jonbodner/channel-talk-2/mutex"
	"github.com/jonbodner/channel-talk-2/mutex"
)

type Val int

func (v Val) read() int {
	return int(v)
}

func (v *Val) inc() {
	*v = *v + 1
}

func TestChannelMutex(t *testing.T) {
	var v Val
	cm := mutex.New()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				cm.Lock()
				vi := v.read()
				runtime.Gosched()
				v.inc()
				runtime.Gosched()
				vp := v.read()
				runtime.Gosched()
				//fmt.Println(vi, vp)
				if vp-vi != 1 {
					t.Errorf("Expected one apart, got %d, %d", vi, vp)
				}
				cm.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
