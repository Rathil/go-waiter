package waiter

import (
	"context"
	"sync"
	"time"
)

type Option struct {
	OneOf    func(param chan int8, cancel chan int8)
	Everyone func(param *sync.WaitGroup)
}

// Create option with timeout
func OptionTimeout(timeout time.Duration) *Option {
	return &Option{
		OneOf: func(param chan int8, cancel chan int8) {
			select {
			case <-time.After(timeout):
				param <- 1
			case <-cancel:
			}
		},
		Everyone: func(param *sync.WaitGroup) {
			param.Add(1)
			go func() {
				select {
				case <-time.After(timeout):
				}
				param.Done()
			}()
		},
	}
}

// Create option with context
func OptionContext(c context.Context) *Option {
	return &Option{
		OneOf: func(param chan int8, cancel chan int8) {
			select {
			case <-c.Done():
				param <- 1
			case <-cancel:
			}
		},
		Everyone: func(param *sync.WaitGroup) {
			param.Add(1)
			go func() {
				select {
				case <-c.Done():
				}
				param.Done()
			}()
		},
	}
}

// Create option with chan
func OptionChan(ch chan int8) *Option {
	return &Option{
		OneOf: func(param chan int8, cancel chan int8) {
			select {
			case <-ch:
				param <- 1
			case <-cancel:
			}
		},
		Everyone: func(param *sync.WaitGroup) {
			param.Add(1)
			go func() {
				select {
				case <-ch:
				}
				param.Done()
			}()
		},
	}
}

// Create option with wait group
func OptionWaitGroup(wg *sync.WaitGroup) *Option {
	return &Option{
		OneOf: func(param chan int8, cancel chan int8) {
			wg.Wait()
			param <- 1
		},
		Everyone: func(param *sync.WaitGroup) {
			param.Add(1)
			go func() {
				wg.Wait()
				param.Done()
			}()
		},
	}
}
