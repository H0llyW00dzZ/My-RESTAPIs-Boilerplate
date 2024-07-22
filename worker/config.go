// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

// NewDoWorkOption defines a functional option for configuring the worker pool.
type NewDoWorkOption[T any] func(*Pool[T])

// WithNumWorkers sets the number of workers in the pool.
//
// Example Usage:
//
//	pool := worker.NewDoWork[uint32](worker.WithNumWorkers(10))
func WithNumWorkers[T any](numWorkers int) NewDoWorkOption[T] {
	return func(wp *Pool[T]) { wp.numWorkers = numWorkers }
}
