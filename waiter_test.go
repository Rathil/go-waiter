package waiter

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestOneOfTimeout(t *testing.T) {
	wait := NewWaiter(
		OptionTimeout(5*time.Second),
		OptionTimeout(10*time.Millisecond),
		OptionTimeout(2*time.Second),
	)
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 10*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 15*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestEveryoneTimeout(t *testing.T) {
	wait := NewWaiter(
		OptionTimeout(15*time.Millisecond),
		OptionTimeout(1*time.Millisecond),
		OptionTimeout(10*time.Millisecond),
	)
	timeStart := time.Now()
	wait.WaitEveryone()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 15*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
}

func TestOneOfContext1(t *testing.T) {
	c, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	wait := NewWaiter(
		OptionContext(c),
	)
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 10*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 15*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestOneOfContext2(t *testing.T) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	wait := NewWaiter(
		OptionContext(c),
	)
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 5*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 9*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestOneOfContextAndTimeout(t *testing.T) {
	c, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	wait := NewWaiter(
		OptionContext(c),
		OptionTimeout(5*time.Millisecond),
	)
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 5*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 9*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestEveryoneContextAndTimeout(t *testing.T) {
	c, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	wait := NewWaiter(
		OptionContext(c),
		OptionTimeout(5*time.Millisecond),
	)
	timeStart := time.Now()
	wait.WaitEveryone()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 10*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
}

func TestOneOfChan1(t *testing.T) {
	ch := make(chan int8)
	wait := NewWaiter(
		OptionChan(ch),
		OptionTimeout(5*time.Millisecond),
	)
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch <- 8
	}()
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 5*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 9*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestOneOfChan2(t *testing.T) {
	ch := make(chan int8)
	wait := NewWaiter(
		OptionChan(ch),
		OptionTimeout(10*time.Millisecond),
	)
	go func() {
		time.Sleep(5 * time.Millisecond)
		ch <- 8
	}()
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 5*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 9*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestEveryoneChan(t *testing.T) {
	ch := make(chan int8)
	wait := NewWaiter(
		OptionChan(ch),
		OptionTimeout(10*time.Millisecond),
	)
	go func() {
		time.Sleep(5 * time.Millisecond)
		ch <- 8
	}()
	timeStart := time.Now()
	wait.WaitEveryone()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 10*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
}

func TestOneOfWG1(t *testing.T) {
	wg := &sync.WaitGroup{}
	wait := NewWaiter(
		OptionWaitGroup(wg),
	)
	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Millisecond)
		wg.Done()
	}()
	timeStart := time.Now()
	wait.WaitOneOf()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 5*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
	if timeEnd.Sub(timeStart) > 9*time.Millisecond {
		t.Fatal("Longer than the specified interval")
	}
}

func TestEveryoneAll(t *testing.T) {
	ch := make(chan int8)
	wg := &sync.WaitGroup{}
	c, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	wait := NewWaiter(
		OptionTimeout(8*time.Millisecond),
		OptionChan(ch),
		OptionWaitGroup(wg),
		OptionContext(c),
	)
	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Millisecond)
		wg.Done()
	}()
	go func() {
		time.Sleep(6 * time.Millisecond)
		ch <- 8
	}()
	timeStart := time.Now()
	wait.WaitEveryone()
	timeEnd := time.Now()
	if timeEnd.Sub(timeStart) < 10*time.Millisecond {
		t.Fatal("Less than the specified interval")
	}
}
