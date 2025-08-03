package task2

import (
	"fmt"
	"sync"
	"time"
)

type Task func()

type TaskResult struct {
	ID        int
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
}

type Scheduler struct {
	tasks     []Task
	results   []TaskResult
	wg        sync.WaitGroup
	startTime time.Time
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks:   make([]Task, 0),
		results: make([]TaskResult, 0),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Run() {
	s.startTime = time.Now()
	s.wg.Add(len(s.tasks))

	for i, task := range s.tasks {
		go func(id int, t Task) {
			defer s.wg.Done()
			start := time.Now()

			t() // 执行任务

			end := time.Now()
			s.results = append(s.results, TaskResult{
				ID:        id,
				Duration:  end.Sub(start),
				StartTime: start,
				EndTime:   end,
			})
		}(i, task)
	}

	s.wg.Wait()
}

func (s *Scheduler) PrintStats() {
	fmt.Println("\n任务执行统计:")
	fmt.Printf("总任务数: %d\n", len(s.tasks))
	fmt.Printf("总执行时间: %v\n", time.Since(s.startTime))

	for _, result := range s.results {
		fmt.Printf("任务%d: 耗时%v (开始: %v 结束: %v)\n",
			result.ID+1,
			result.Duration,
			result.StartTime.Format("15:04:05.000"),
			result.EndTime.Format("15:04:05.000"))
	}
}

/*
func main() {
	scheduler := NewScheduler()

	// 添加示例任务
	scheduler.AddTask(func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("任务1完成")
	})

	scheduler.AddTask(func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("任务2完成")
	})

	scheduler.AddTask(func() {
		time.Sleep(700 * time.Millisecond)
		fmt.Println("任务3完成")
	})

	fmt.Println("开始执行任务...")
	scheduler.Run()
	scheduler.PrintStats()
}
*/
