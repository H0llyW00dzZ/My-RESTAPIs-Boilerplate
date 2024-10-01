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
	numWorkers int            // Store the number of workers
	activeJobs int32          // Track the number of active jobs
	isRunning  uint32
	mu         sync.Mutex
	// Store registered job functions
	//
	// Note: this optional it can bound to other instead of [fiber.Ctx] (e.g, database for streaming html hahaha).
	//
	// TODO: Improve this for more flexibility to allow other types instead of currently only supporting "func(*fiber.Ctx)"
	registeredJobs map[string]func(*fiber.Ctx) Job[T]

	// Channel options
	jobChannelOpts    []ChanOption[Job[T]]
	resultChannelOpts []ChanOption[T]
	errorChannelOpts  []ChanOption[error]
	// Configurable idle check interval
	idleCheckInterval time.Duration
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
func NewDoWork[T any](opts ...NewDoWorkOption[T]) *Pool[T] {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &Pool[T]{
		ctx:            ctx,
		cancel:         cancel,
		wg:             sync.WaitGroup{},
		numWorkers:     NumWorkers,
		activeJobs:     0,
		isRunning:      0,
		mu:             sync.Mutex{},
		registeredJobs: make(map[string]func(*fiber.Ctx) Job[T]),

		// Initialize channels with default options
		jobs:              make(chan Job[T], NumWorkers),
		results:           make(chan T, NumWorkers),
		errors:            make(chan error, NumWorkers),
		idleCheckInterval: DefaultWorkerSleepTime, // Default value
	}

	// Apply functional options
	for _, opt := range opts {
		opt(wp)
	}

	// Apply channel options
	wp.applyChanOptions()

	return wp
}

// Stop gracefully shuts down the worker pool
func (wp *Pool[T]) Stop() {
	if wp.atomicStop() {
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
	if !wp.IsRunning() {
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
	if !wp.atomicStart() {
		return
	}
	// Note: this used std logger, due it not possible import internal package in the backend to outside (not allowed).
	log.Print("Worker pool started.")
	go func() {
		defer func() {
			atomic.StoreUint32(&wp.isRunning, 0)
			log.Print("Worker pool exiting.")
		}()

		// Use the WaitGroup to wait for workers to start
		wp.wg.Add(wp.numWorkers)
		for w := 0; w < wp.numWorkers; w++ {
			go func() {
				defer wp.wg.Done() // Signal when a worker is ready
				for {
					select {
					case job, ok := <-wp.jobs:
						// Check if the channel is closed
						if !ok {
							return // Exit the worker goroutine
						}
						result, err := job.Execute(wp.ctx)
						if err != nil {
							log.Printf("Error executing job: %v", err)

							// Safe error sending: check if the pool is still running
							if wp.IsRunning() {
								wp.errors <- err // Signal error
							}
						} else {
							// Safe result sending
							if wp.IsRunning() {
								wp.results <- result
								log.Printf("worker finished job with result: %v", result)
							}
						}
					case <-wp.ctx.Done(): // Listen for context cancellation for shutdown
						return
					}
				}
			}()
		}

		// Idle worker monitoring and shutdown logic SHOULD BE HERE!
		// Wait for all workers to signal they are ready
		wp.wg.Wait() // <- This Correct reallocation for long-running (e.g, zer0-downtime, till next billion years, handling http traffic "keep-alive") task.

		idleTicker := time.NewTicker(wp.idleCheckInterval) // Use default or configured
		defer idleTicker.Stop()

		for {
			select { // Check for idleness every second if use default or configured
			case <-idleTicker.C:
				if !wp.IsRunning() {
					wp.Stop()
					return // Exit the loop when the pool is stopped
				}
			}
			// Now we wait for all workers to be done before checking if
			// we need to shut down
			// Note: This will be problematic (e.g, crash cause the caller (other goroutine) sending to closed channel) for long running-task for example when the worker pool has shutdown
			// there is no chance for the caller (other goroutine) to submit/start again when implement the execute function,
			// unless application must be restart (e.g, shutdown then start again), so disabled by commented out here.
			//wp.wg.Wait() // <- Disabled due it problematic for long-running task.
		}
	}()
}

// IsRunning checks if the worker pool is currently running.
//
// It returns true if the pool is running, false otherwise.
func (wp *Pool[T]) IsRunning() bool {
	return atomic.LoadUint32(&wp.isRunning) == 1
}

// applyChanOptions applies the configured channel options.
func (wp *Pool[T]) applyChanOptions() {
	for _, opt := range wp.jobChannelOpts {
		wp.applyChanOption(wp.jobs, opt)
	}
	for _, opt := range wp.resultChannelOpts {
		wp.applyChanOption(wp.results, opt)
	}
	for _, opt := range wp.errorChannelOpts {
		wp.applyChanOption(wp.errors, opt)
	}
}

// applyChanOption applies a channel option to a channel.
//
// Note: This not using pointer that's why it implement this function
func (wp *Pool[T]) applyChanOption(ch any, opt any) {
	switch o := opt.(type) {
	case ChanOption[Job[T]]:
		o(ch.(chan Job[T]))
	case ChanOption[T]:
		o(ch.(chan T))
	case ChanOption[error]:
		o(ch.(chan error))
	default:
		panic(fmt.Sprintf("unsupported channel option type: %T", opt))
	}
}

// atomicStart atomically sets the isRunning flag to 1 if it is currently 0.
// It returns true if the operation was successful (i.e., isRunning was 0 and is now 1),
// indicating that the worker pool has started running.
// It returns false if the worker pool is already running (i.e., isRunning was already 1).
func (wp *Pool[T]) atomicStart() bool {
	return atomic.CompareAndSwapUint32(&wp.isRunning, 0, 1)
}

// atomicStop atomically sets the isRunning flag to 0 if it is currently 1.
// It returns true if the operation was successful (i.e., isRunning was 1 and is now 0),
// indicating that the worker pool has stopped running.
// It returns false if the worker pool is already stopped (i.e., isRunning was already 0).
func (wp *Pool[T]) atomicStop() bool {
	return atomic.CompareAndSwapUint32(&wp.isRunning, 1, 0)
}
