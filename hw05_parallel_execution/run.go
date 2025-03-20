package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Если m <= 0 – ошибки игнорируются.
	shouldLimitErrors := m > 0
	var errorsCount int64

	taskCh := make(chan Task)

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil && shouldLimitErrors {
					atomic.AddInt64(&errorsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if shouldLimitErrors && atomic.LoadInt64(&errorsCount) >= int64(m) {
			break
		}
		taskCh <- task
	}

	close(taskCh)

	wg.Wait()

	if shouldLimitErrors && atomic.LoadInt64(&errorsCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
