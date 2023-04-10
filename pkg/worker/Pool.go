package worker

import (
	"TrackMaster/pkg"
)

const (
	MaxWorkers = 10  // 最大goroutine数
	MaxQueue   = 100 // 最大队列长度
)

type Job interface {
	Execute() *pkg.Error
}

// Pool 定义用于执行任务的goroutine池
type Pool struct {
	MaxWorkers int               // 最大goroutine数
	MaxQueue   int               // 最大队列长度
	Jobs       chan Job          // 新任务通道
	Errors     chan<- *pkg.Error // 错误通道
}

func NewWorkerPool(errCh chan<- *pkg.Error) *Pool {
	return &Pool{
		MaxWorkers: MaxWorkers,
		MaxQueue:   MaxQueue,
		Jobs:       make(chan Job, MaxQueue),
		Errors:     errCh,
	}
}

func (wp *Pool) Begin() {
	for i := 0; i < wp.MaxWorkers; i++ {
		go func() {
			for job := range wp.Jobs {
				if err := job.Execute(); err != nil {
					wp.Errors <- err
				}
			}
		}()
	}
}
