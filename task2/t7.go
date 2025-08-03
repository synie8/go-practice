package task2

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type AtomicCounter struct {
	count int64
}

func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.count, 1)
}
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.count)
}
func Exe() {
	var wg sync.WaitGroup
	atom := AtomicCounter{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				atom.Inc()
			}
			time.Sleep(time.Millisecond * 10)
		}()

	}
	wg.Wait()
	fmt.Println("final value -> ", atom.Value())
}
