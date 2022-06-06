package util

import "sync"

type Nursery struct {
	waitGroup *sync.WaitGroup
}

func NewNursery() *Nursery {
	return &Nursery{waitGroup: &sync.WaitGroup{}}
}

func (n *Nursery) Start(f func()) {
	n.waitGroup.Add(1)
	go func() {
		f()
		n.waitGroup.Done()
	}()
}

func (n *Nursery) Wait() {
	n.waitGroup.Wait()
}
