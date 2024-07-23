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
//	pool := worker.NewDoWork(worker.WithNumWorkers[uint32](10))
func WithNumWorkers[T any](numWorkers int) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.numWorkers = numWorkers
	}
}

// ChanOption defines a functional option for configuring channels.
type ChanOption[C any] func(ch chan C)

// WithJobChannelOptions configures the job channel.
func WithJobChannelOptions[T any](opts ...ChanOption[Job[T]]) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.jobChannelOpts = opts
	}
}

// WithResultChannelOptions configures the results channel.
func WithResultChannelOptions[T any](opts ...ChanOption[T]) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.resultChannelOpts = opts
	}
}

// WithErrorChannelOptions configures the errors channel.
func WithErrorChannelOptions[T any](opts ...ChanOption[error]) NewDoWorkOption[T] {
	return func(wp *Pool[T]) {
		wp.errorChannelOpts = opts
	}
}

// WithChanBuffer sets the buffer size of a channel.
func WithChanBuffer[C any](bufferSize int) ChanOption[C] {
	return func(ch chan C) {
		ch = make(chan C, bufferSize)
	}
}
