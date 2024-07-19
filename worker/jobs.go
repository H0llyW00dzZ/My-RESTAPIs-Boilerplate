// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	// ErrFailedToGetSomething is returned when failed to get something..
	ErrFailedToGetSomething = errors.New("worker failed to get something from job")
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

// job represents a unit of work for the worker pool.
type job struct {
	c *fiber.Ctx // this optional it can bound to other (e.g, database for streaming html hahaha).
}
