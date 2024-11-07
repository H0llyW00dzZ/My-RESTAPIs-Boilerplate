// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"errors"
	"time"
)

var (
	// ErrFailedToGetSomething is returned when failed to get something..
	ErrFailedToGetSomething = errors.New("worker: failed to get something from job")
	// ErrJobsNotFound is returned when a job with the specified name is not registered with the worker pool.
	ErrJobsNotFound = errors.New("worker: job not found")
	// ErrorInvalidJobType is returned when a job function returns a type that is not expected or supported by the worker.
	ErrorInvalidJobType = errors.New("worker: invalid job function return type")
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
	//
	// Additionally, for Horizontal Pod Autoscaling (HPA) in Kubernetes, it's better to set
	// NumWorkers to 300 and adjust the CPU to 350m ~ 450m. This allows the application to scale
	// to around 5 pods, ensuring effective synchronization, unlike stateful applications, which
	// cannot scale as easily as stateless ones.
	// Stateful applications are typically used for databases, and using them for web services
	// is not recommended (bad).
	NumWorkers = 1
)

// Default Worker Configuration
const (
	// Note: This Recommended and Suitable for handling long traffic
	// (e.g, long request http till next billion years then it stop),
	// high traffic (e.g, many request incoming from different client source, unlike long traffic),
	// a perfect scheduler (e.g, 24/7 automated updated value for example in database), other worker (e.g, background, assistant garbage collector).
	DefaultWorkerSleepTime = 1 * time.Minute
)

// NewDoWorkOption defines a functional option for configuring the worker pool.
type NewDoWorkOption[T any] func(*Pool[T])

// WithNumWorkers sets the number of workers in the pool.
//
// Example Usage:
//
//	pool := worker.NewDoWork(worker.WithNumWorkers[uint32](10))
//
// Note: It's important to note that increasing the number of workers and buffer sizes can potentially improve the performance and concurrency of the worker pool,
// but it also depends on the specific use case and the available system resources.
// It's recommended to tune these values based on go application's requirements and performance characteristics.
func WithNumWorkers[T any](numWorkers int) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.numWorkers = numWorkers
	}
}

// ChanOption defines a functional option for configuring channels.
type ChanOption[C any] func(ch chan C)

// WithJobChannelOptions configures the job channel.
//
// Note: It's important to note that increasing the number of workers and buffer sizes can potentially improve the performance and concurrency of the worker pool,
// but it also depends on the specific use case and the available system resources.
// It's recommended to tune these values based on go application's requirements and performance characteristics.
func WithJobChannelOptions[T any](opts ...ChanOption[Job[T]]) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.jobChannelOpts = opts
	}
}

// WithResultChannelOptions configures the results channel.
//
// Note: It's important to note that increasing the number of workers and buffer sizes can potentially improve the performance and concurrency of the worker pool,
// but it also depends on the specific use case and the available system resources.
// It's recommended to tune these values based on go application's requirements and performance characteristics.
func WithResultChannelOptions[T any](opts ...ChanOption[T]) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.resultChannelOpts = opts
	}
}

// WithErrorChannelOptions configures the errors channel.
//
// Note: It's important to note that increasing the number of workers and buffer sizes can potentially improve the performance and concurrency of the worker pool,
// but it also depends on the specific use case and the available system resources.
// It's recommended to tune these values based on go application's requirements and performance characteristics.
func WithErrorChannelOptions[T any](opts ...ChanOption[error]) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.errorChannelOpts = opts
	}
}

// WithChanBuffer sets the buffer size of a channel.
//
// Note: It's important to note that increasing the number of workers and buffer sizes can potentially improve the performance and concurrency of the worker pool, but it also depends on the specific use case and the available system resources.
// It's recommended to tune these values based on your application's requirements and performance characteristics.
func WithChanBuffer[C any](bufferSize int) ChanOption[C] {
	return func(ch chan C) {
		ch = make(chan C, bufferSize)
	}
}

// WithIdleCheckInterval sets the interval at which the worker pool checks
// for idleness and potentially shuts down.
//
// Note: It's important to note that increasing the number of workers and buffer sizes can potentially improve the performance and concurrency of the worker pool, but it also depends on the specific use case and the available system resources.
// It's recommended to tune these values based on your application's requirements and performance characteristics.
func WithIdleCheckInterval[T any](interval time.Duration) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.idleCheckInterval = interval
	}
}
