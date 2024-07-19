// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	// ErrFailedToGetSomething is returned when failed to get something..
	ErrFailedToGetSomething = errors.New("worker failed to get something from job")
	// ErrJobsNotFound is returned when a job with the specified name is not registered with the worker pool.
	ErrJobsNotFound = errors.New("job not found")
)

const (
	// NumWorkers it for set how many worker, for example I am using 100 worker
	// that used for handle high traffic + large go application (600+ files) not waste memory.
	//
	// Also note that there is the price:
	//
	// 100 worker = 100mb ~ 150mb++ memory consumed (Approx)
	//
	// under 50 worker still consider cheap.
	NumWorkers = 5
)

// Default Worker Configuration
const (
	// Note: This Recommended and Suitable for handling long traffic
	// (e.g, long request http till next billion years then it stop),
	// high traffic (e.g, many request incoming from different client source, unlike long traffic),
	// a perfect scheduler (e.g, 24/7 automated updated value for example in database), other worker (e.g, background, assistant garbage collector).
	DefaultWorkerSleepTime = 1 * time.Second
)

// Job represents a unit of work for the worker pool.
type Job interface {
	// Execute runs the job, returning a result (or an error if it failed)
	Execute(ctx context.Context) (string, error)
}

// RegisterJob adds a new job function to the pool.
//
// Example:
//
//	pool.RegisterJob("myStreamingJob", func(c *fiber.Ctx) worker.Job {
//	    return &MyStreamingJob{c: c}
//	})
//
// Example with Init (Recommended when put in somewhere e.g, outside of worker package):
//
//	func init() {
//		pool.RegisterJob("myStreamingJob", func(c *fiber.Ctx) worker.Job {
//			return &MyStreamingJob{c: c}
//		})
//	}
//
// Execute the job:
//
//	func (s *MyStreamingJob) Execute(ctx context.Context) (string, error) {
//		// Your Function Poggers...
//		return "", worker.ErrFailedToGetSomething
//	}
//
// Then call the worker.Submit see (worker.NewDoWork).
//
// Note: New design it more flexibility, unlike previous design.
func (wp *Pool) RegisterJob(name string, jobFunc func(*fiber.Ctx) Job) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.registeredJobs[name] = jobFunc
}
