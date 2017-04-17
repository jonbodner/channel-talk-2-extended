package mutex

type Mutex interface {
	Lock()
	Unlock()
}

type channelMutex chan bool

func New() Mutex {
	cm := make(channelMutex, 1)
	cm <- true
	return cm
}

func (cm channelMutex) Lock() {
	<-cm
}

func (cm channelMutex) Unlock() {
	cm <- true
}

