package worker

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"taskflow/internal/application"
	"taskflow/internal/domain"
	"time"
)

type Result struct {
	TaskID  domain.TaskID
	Success bool
	Err     error
}

type Pool struct {
	workers     int
	jobChan     chan *domain.Task
	resultChan  chan Result
	service     *application.TaskService
	collectorWg sync.WaitGroup
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	stopOnce    sync.Once
}

func NewPool(workers int, service *application.TaskService) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workers:    workers,
		jobChan:    make(chan *domain.Task, workers),
		resultChan: make(chan Result, workers*2),
		service:    service,
		ctx:        ctx,
		cancel:     cancel,
		stopOnce:   sync.Once{},
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		// adds one to the WaitGroup counter for each worker goroutine
		p.wg.Add(1)
		go p.worker()
	}
	p.collectorWg.Add(1)
	go p.collector()
	go p.monitor()
}

func (p *Pool) worker() {
	defer p.wg.Done()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for task := range p.jobChan {
		startTime := time.Now()
		fmt.Printf("Engine: Starting task %s (Priority: %s)\n", task.ID, task.Priority)
		// Random sleep between 1 and 5 seconds
		sleepTime := time.Duration(r.Intn(5)+1) * time.Second
		select {
		case <-p.ctx.Done():
			p.resultChan <- Result{TaskID: task.ID, Success: false, Err: p.ctx.Err()}
		case <-time.After(sleepTime):
			success := rand.Float64() < 0.8 // 80% chance of success
			fmt.Printf("Engine: Task %s processing finished.\n", task.ID)
			var err error
			if success {
				err = p.service.MarkDone(p.ctx, task.ID)
			} else {
				err = p.service.MarkFailed(p.ctx, task.ID)
			}
			p.resultChan <- Result{TaskID: task.ID, Success: success, Err: err}
		}
		endTime := time.Now()
		fmt.Printf("Engine: Task %s completed in %v\n", task.ID, endTime.Sub(startTime))
	}
}

func (p *Pool) monitor() {
	p.wg.Wait()
	close(p.resultChan)
}

func (p *Pool) collector() {
	defer p.collectorWg.Done()

	for result := range p.resultChan {
		if result.Err != nil {
			fmt.Printf("[ERROR] Task %s failed: %v\n", result.TaskID, result.Err)
		} else if result.Success {
			fmt.Printf("[DONE]  Task %s completed successfully\n", result.TaskID)
		} else {
			fmt.Printf("[FAIL]  Task %s was marked as failed\n", result.TaskID)
		}
	}
}

func (p *Pool) Submit(task *domain.Task) error {
	select {
	case p.jobChan <- task:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("pool is stopped, cannot submit task %s", task.ID)
	}
}

func (p *Pool) Stop() {
	p.stopOnce.Do(func() {
		p.cancel()
		close(p.jobChan)
		p.wg.Wait()

		p.collectorWg.Wait()
	})
}
