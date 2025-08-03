package task2

import (
	"fmt"
	"sync"
)

func Send(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
		fmt.Println("发送的数据****", i)

	}

	close(ch)
}

func Receive(ch <-chan int) {
	for v := range ch {
		fmt.Println("接受到的数据****", v)

	}
}

func Producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 100; i++ {
		ch <- i
		fmt.Printf("生产者发送: %d\n", i)
	}
	close(ch)
}

func Consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("消费者接收: %d\n", num)
	}
}

/*
func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 10) // 创建缓冲大小为10的通道

	wg.Add(2)
	go producer(ch, &wg)
	go consumer(ch, &wg)

	wg.Wait()
	fmt.Println("程序结束")
}

*/

/*


func main() {
	ch := make(chan int, 10)

	go task2.Send(ch)
	go task2.Receive(ch)

	timeOut := time.After(time.Second * 3)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				fmt.Println("ch已关闭")
				return
			}
		case <-timeOut:
			fmt.Println("操作超时")
			return
		default:
			fmt.Println("等待中*******")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

*/
