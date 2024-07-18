// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Pool manages a pool of goroutines for work.
type Pool struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup // Use a single WaitGroup for both startup & shutdown
	jobs       chan job
	results    chan string // TODO: Improve this, instead of string.
	activeJobs int32       // Track the number of active jobs
	isRunning  uint32
	mu         sync.Mutex
}

// NewDoWork creates a new pool and do work just like human being.
//
// Note: This required global variable for example put this in somewhere:
//
//	var pool = worker.NewDoWork()
//
// also note that this safe and idiom go.
func NewDoWork() *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &Pool{
		ctx:        ctx,
		cancel:     cancel,
		wg:         sync.WaitGroup{},
		jobs:       make(chan job, NumWorkers),
		results:    make(chan string, NumWorkers),
		activeJobs: 0,
		isRunning:  0,
		mu:         sync.Mutex{},
	}
	return wp
}

// Stop gracefully shuts down the worker pool
func (wp *Pool) Stop() {
	if atomic.CompareAndSwapUint32(&wp.isRunning, 1, 0) {
		wp.mu.Lock()
		defer wp.mu.Unlock()
		log.Print("Shutting down worker pool...")
		close(wp.jobs) // Signal workers to stop
		wp.cancel()    // Cancel the context
		wp.wg.Wait()   // Wait for workers to finish
		log.Print("Worker pool shut down.")
	}
}

// Submit a job to the worker pool
func (wp *Pool) Submit(c *fiber.Ctx) (string, error) {
	if atomic.LoadUint32(&wp.isRunning) == 0 {
		wp.Start()
	}

	atomic.AddInt32(&wp.activeJobs, 1)        // Increment job counter
	defer atomic.AddInt32(&wp.activeJobs, -1) // Decrement on function exit

	wp.jobs <- job{c: c}
	something := <-wp.results
	if something == "" {
		return "", ErrFailedToGetSomething
	}
	return something, nil
}

// Start a job to the worker pool
func (wp *Pool) Start() {
	if !atomic.CompareAndSwapUint32(&wp.isRunning, 0, 1) {
		return
	}
	log.Print("Worker pool started.")
	go func() {
		defer atomic.StoreUint32(&wp.isRunning, 0)
		defer log.Print("Worker pool exiting.")

		// Use the WaitGroup to wait for workers to start
		wp.wg.Add(NumWorkers)
		for w := 0; w < NumWorkers; w++ {
			go func() {
				defer wp.wg.Done() // Signal when a worker is ready
				// Important: Put Function here so they will starting doing work just like human being
			}()
		}

		// Wait for all workers to signal they are ready
		wp.wg.Wait()

		for {
			time.Sleep(1 * time.Second) // Check for idleness every second
			if atomic.LoadInt32(&wp.activeJobs) == 0 {
				wp.Stop()
				return // Exit the loop when the pool is stopped
			}
		}
	}()
}
