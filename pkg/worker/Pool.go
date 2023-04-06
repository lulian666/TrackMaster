package worker

import (
	"TrackMaster/model"
	"TrackMaster/pkg/track"
	"fmt"
)

const (
	MaxWorkers = 10  // 最大goroutine数
	MaxQueue   = 100 // 最大队列长度
)

type Duty string

const (
	Fetch Duty = "Fetch"
	Test  Duty = "Test"
)

// Task 这个task指的是什么？
type Task struct {
	ID     string // 任务id
	Cancel chan struct{}

	Type   Duty
	Record model.Record
}

// Pool 定义用于执行任务的goroutine池
type Pool struct {
	MaxWorkers int          // 最大goroutine数
	MaxQueue   int          // 最大队列长度
	Tasks      chan *Task   // 任务通道
	Errors     chan<- error // 错误通道
}

func NewWorkerPool(errCh chan<- error) *Pool {
	return &Pool{
		MaxWorkers: MaxWorkers,
		MaxQueue:   MaxQueue,
		Tasks:      make(chan *Task, MaxQueue),
		Errors:     errCh,
	}
}

func (wp *Pool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		go func() {
			for task := range wp.Tasks {
				if err := wp.DoTask(task); err != nil {
					wp.Errors <- err
				}
			}
		}()
	}
}

func (wp *Pool) DoTask(task *Task) error {
	if task.Type == Fetch {
		count, err := track.FetchNewLog(task.Record)
		if err != nil {
			task.Cancel <- struct{}{} // 如果读log出错，这个task就取消 todo 在哪里取消
			return err
		}

		if count > 0 {
			testTask := &Task{
				ID:     task.ID,
				Cancel: make(chan struct{}),
				Type:   Test,
				Record: task.Record,
			}
			select {
			case wp.Tasks <- testTask:
				fmt.Println("Test Task has been added to the queue")
			default:
				fmt.Println("Too many Tasks")
			}
		}
	}

	if task.Type == Test {
		err := track.TestNewLog(task.Record)
		if err != nil {
			task.Cancel <- struct{}{} // 如果读log出错，这个task就取消
			return err
		}
	}

	return nil
}
