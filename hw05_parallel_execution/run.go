package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	tasksChan := make(chan Task, len(tasks)) // канал для списка задач
	done := make(chan struct{})              // канал сигнальный для закрытия всех горутин
	errorsCount := 0                         // подсчет кол-ва ошибок в задачах
	stop := false                            // флаг остановки горутин при достижении кол-ва ошибок errorsCount числа m

	// Подаем задачи в канал
	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	// Запускаем n воркеров
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(tasksChan, done, &wg, &errorsCount, &stop, mu, m)
	}

	wg.Wait()

	if errorsCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(
	tasksChan <-chan Task,
	done chan struct{},
	wg *sync.WaitGroup,
	errorsCount *int,
	stop *bool,
	mu *sync.Mutex,
	m int,
) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		case task, ok := <-tasksChan:
			if !ok {
				return
			}
			if err := task(); err != nil {
				mu.Lock()
				*errorsCount++
				if *errorsCount >= m {
					if !*stop {
						*stop = true
						close(done)
					}
					mu.Unlock()
					return
				}
				mu.Unlock()
			}
		}
	}
}
