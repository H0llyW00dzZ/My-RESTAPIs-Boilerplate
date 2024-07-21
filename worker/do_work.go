// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Pool manages a pool of goroutines for work.
type Pool[T any] struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup // Use a single WaitGroup for both startup & shutdown
	jobs       chan Job[T]    // Queue for jobs
	results    chan T         // Results channel, now generic, it more easier instead of only string.
	errors     chan error     // Error channel collections, each worker had their own error channel for communication same as results channel (e.g, 1000 worker/goroutines then 1000 error).
	activeJobs int32          // Track the number of active jobs
	isRunning  uint32
	mu         sync.Mutex
	// Store registered job functions
	//
	// Note: this optional it can bound to other instead of [fiber.Ctx] (e.g, database for streaming html hahaha).
	registeredJobs map[string]func(*fiber.Ctx) Job[T]
}

// NewDoWork creates a new pool and do work just like human being.
//
// Note: This required global variable for example put this in somewhere:
//
//	var pool = worker.NewDoWork[string]()
//
// Then Call it Example:
//
//	func myWorkerDoingStreaming(c *fiber.Ctx) error {
//		streamingHTML, err := pool.Submit(c)
//		  if err != nil {
//			 // handle error you poggers
//		  }
//
//		 // Send the response
//		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
//		// Send the rendered HTML content as a response with the appropriate status code.
//		return c.Status(statusCode).SendString(buf.String())
//	}
//
// Also note that this safe and idiom go.
func NewDoWork[T any]() *Pool[T] {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &Pool[T]{
		ctx:            ctx,
		cancel:         cancel,
		wg:             sync.WaitGroup{},
		jobs:           make(chan Job[T], NumWorkers),
		results:        make(chan T, NumWorkers),
		errors:         make(chan error, NumWorkers),
		activeJobs:     0,
		isRunning:      0,
		mu:             sync.Mutex{},
		registeredJobs: make(map[string]func(*fiber.Ctx) Job[T]),
	}
	return wp
}

// Stop gracefully shuts down the worker pool
func (wp *Pool[T]) Stop() {
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
func (wp *Pool[T]) Submit(c *fiber.Ctx, jobName string) (T, error) {
	if atomic.LoadUint32(&wp.isRunning) == 0 {
		wp.Start()
	}

	wp.mu.Lock()
	jobFunc, ok := wp.registeredJobs[jobName]
	wp.mu.Unlock()

	if !ok {
		var zero T // Might want to return an appropriate "zero" value for generic type here.
		return zero, fmt.Errorf("%w: %s", ErrJobsNotFound, jobName)
	}

	atomic.AddInt32(&wp.activeJobs, 1)        // Increment job counter
	defer atomic.AddInt32(&wp.activeJobs, -1) // Decrement on function exit

	job := jobFunc(c)
	wp.jobs <- job
	select {
	case result := <-wp.results:
		return result, nil
	case err := <-wp.errors:
		return *new(T), err // Return the zero value of T and the error
	}
}

// Start a job to the worker pool
func (wp *Pool[T]) Start() {
	if !atomic.CompareAndSwapUint32(&wp.isRunning, 0, 1) {
		return
	}
	// Note: this used std logger, due it not possible import internal package in the backend to outside (not allowed).
	log.Print("Worker pool started.")
	go func() {
		defer atomic.StoreUint32(&wp.isRunning, 0)
		defer log.Print("Worker pool exiting.")

		// Use the WaitGroup to wait for workers to start
		wp.wg.Add(NumWorkers)
		for w := 0; w < NumWorkers; w++ {
			go func() {
				defer wp.wg.Done() // Signal when a worker is ready
				for job := range wp.jobs {
					result, err := job.Execute(wp.ctx)
					if err != nil {
						log.Printf("Error executing job: %v", err)
						wp.errors <- err // Signal error
					} else {
						wp.results <- result
						log.Printf("worker finished job with result: %v", result)
					}
				}
			}()
		}

		// Wait for all workers to signal they are ready
		wp.wg.Wait()

		for {
			time.Sleep(DefaultWorkerSleepTime) // Check for idleness every second
			if atomic.LoadInt32(&wp.activeJobs) == 0 {
				wp.Stop()
				return // Exit the loop when the pool is stopped
			}
		}
	}()
}

// IsRunning checks if the worker pool is currently running.
//
// It returns true if the pool is running, false otherwise.
func (wp *Pool[T]) IsRunning() bool {
	return atomic.LoadUint32(&wp.isRunning) == 1
}
