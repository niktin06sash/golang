package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID    int
	Value int
}
type Result struct {
	TaskID    int
	TaskValue int
}
type WorkerPool struct {
	Countworker int
	Taskchan    chan Task
	Resultchan  chan Result
	wg          sync.WaitGroup
	Context     context.Context
	Cancel      context.CancelFunc
}

func NewWorkerPool(numworkers int, contex context.Context) *WorkerPool {
	ctx, cancel := context.WithCancel(contex)
	return &WorkerPool{
		Countworker: numworkers,
		Taskchan:    make(chan Task, 100),
		Resultchan:  make(chan Result, 100),
		Context:     ctx,
		Cancel:      cancel,
	}
}
func (wp *WorkerPool) Run() {
	for i := 1; i <= wp.Countworker; i++ {
		wp.wg.Add(1)
		go wp.DoTask(i)
	}
	go wp.GetResult()
}
func (wp *WorkerPool) Stop() {
	close(wp.Taskchan)
	wp.Cancel()
	wp.wg.Wait()
	close(wp.Resultchan)
}
func (wp *WorkerPool) DoTask(id int) {
	defer wp.wg.Done()
	fmt.Printf("Worker %d: Starting\n", id)
	for {
		select {
		case newtask, ok := <-wp.Taskchan:
			if !ok {
				return
			}
			fmt.Printf("Worker %d: Received task ID=%d, Value=%d\n", id, newtask.ID, newtask.Value)
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond
			select {
			case <-time.After(delay):
				resultValue := newtask.Value * 2
				fmt.Printf("Worker %d: Task ID=%d processed, result=%d\n", id, newtask.ID, resultValue)
				wp.Resultchan <- Result{TaskID: newtask.ID, TaskValue: resultValue}
			case <-wp.Context.Done():
				fmt.Printf("TaskChan closed!")
				return
			}
		case <-wp.Context.Done():
			fmt.Printf("TaskChan closed!")
			return
		}
	}
}
func (wp *WorkerPool) GetResult() {
	for result := range wp.Resultchan {
		fmt.Printf("Result Handler: Task ID=%d, Result=%d\n", result.TaskID, result.TaskValue)
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	workerPool := NewWorkerPool(5, ctx)
	workerPool.Run()
	numbertask := 100
	for i := 1; i <= numbertask; i++ {
		workerPool.Taskchan <- Task{ID: i, Value: rand.Intn(1000)}
	}
	time.Sleep(10 * time.Second)
	go workerPool.Stop()

}
