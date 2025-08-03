package task2

import (
	"sync"
)

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

/*
func main() {
	counter := SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				counter.Inc()
			}
			time.Sleep(time.Millisecond * 10)
		}()
	}
	wg.Wait()
	fmt.Println("final count: ", counter.Value())
}
*/
