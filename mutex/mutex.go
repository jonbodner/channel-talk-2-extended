package mutex

type ChannelMutex chan bool

func New() ChannelMutex {
	cm := make(ChannelMutex, 1)
	cm <- true
	return cm
}

func (cm ChannelMutex) Lock() {
	<-cm
}

func (cm ChannelMutex) Unlock() {
	cm <- true
}

