package mutex_test

import (
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

func (v *Val) increment() {
	*v = *v + 1
}

type fakeMutex struct{}

func (fm *fakeMutex) Lock() {

}

func (fm *fakeMutex) Unlock() {

}

func TestChannelMutexFail(t *testing.T) {
	testIt(&fakeMutex{}, t)
}

func TestChannelMutex(t *testing.T) {
	cm := mutex.New()
	testIt(cm, t)
}

func testIt(cm mutex.Mutex, t *testing.T) {
	var v Val
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				cm.Lock()
				vi := v.read()
				runtime.Gosched()
				v.increment()
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
