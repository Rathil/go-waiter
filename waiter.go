package waiter

import (
	"sync"
)

type Waiter struct {
	opt []*Option
}

// Create new waiter
func NewWaiter(opt ...*Option) *Waiter {
	return &Waiter{
		opt: opt,
	}
}

// Waiting for any event
func (a *Waiter) WaitOneOf() {
	wait := make(chan int8, len(a.opt))
	cancel := make(chan int8, len(a.opt)-1)
	for _, opt := range a.opt {
		go opt.OneOf(wait, cancel)
	}
	defer func() {
		for i := 1; i < len(a.opt); i++ {
			cancel <- 1
		}
	}()
	select {
	case <-wait:
	}
}

// Waiting for all events
func (a *Waiter) WaitEveryone() {
	wg := &sync.WaitGroup{}
	for _, opt := range a.opt {
		opt.Everyone(wg)
	}
	wg.Wait()
}
