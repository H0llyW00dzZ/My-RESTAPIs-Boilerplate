// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Job represents a unit of work for the worker pool.
type Job[T any] interface {
	// Execute runs the job, returning a result (or an error if it failed)
	Execute(ctx context.Context) (T, error)
}

// RegisterJob adds a new job function to the pool.
//
// Example:
//
//	pool.RegisterJob("myStreamingJob", func(c *fiber.Ctx) worker.Job[string] {
//	    return &MyStreamingJob[string]{c: c}
//	})
//
// Example with Init (Recommended when put in somewhere e.g, outside of worker package):
//
//	func init() {
//		pool.RegisterJob("myStreamingJob", func(c *fiber.Ctx) worker.Job[string] {
//			return &MyStreamingJob[string]{c: c}
//		})
//	}
//
// Execute the job:
//
//	func (s *MyStreamingJob[T]) Execute(ctx context.Context) (T, error) {
//
//		// Your Function Poggers...
//	    // Perform the job logic here
//	    // For example, make database queries, process data, or interact with external services.
//
//	    // Return the result of the job and an error if any.
//	    return someResult, nil // Replace someResult with the actual result of your job.
//	}
//
// Then call the worker.Submit see (worker.NewDoWork).
//
// Note: The new design is more flexible (unlike previous design) and eliminates the need for explicit mutex locks/unlocks when implementing the [Execute] function.
// This is because the use of channels and atomic operations in the worker pool ensures that data is accessed and modified safely without the risk of data races.
// Also note that, in the example "Execute the job," it's possible to spawn additional goroutines for a more idiomatic Go approach,
// allowing them to communicate with each other while performing some work.
//
// In general, it's important to remember the following principles when working with shared memory and concurrency in Go:
//   - Don't communicate by sharing memory; share memory by communicating.
//     This means using channels or other synchronization primitives to communicate and exchange data between goroutines, rather than directly accessing shared memory locations.
//   - Use atomic operations to modify shared data safely. Atomic operations guarantee that each access to the shared data is atomic, meaning that the value of the data is always consistent.
//   - Use synchronization primitives such as mutexes or channels to control access to shared resources and prevent data races.
func (wp *Pool[T]) RegisterJob(name string, jobFunc func(*fiber.Ctx) Job[T]) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.registeredJobs[name] = jobFunc
}
