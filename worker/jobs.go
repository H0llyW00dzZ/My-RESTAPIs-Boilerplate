// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	// ErrFailedToGetSomething is returned when failed to get something..
	ErrFailedToGetSomething = errors.New("worker failed to get something from job")
)

const (
	// NumWorkers it for set how many worker, for example I am using 100 worker
	NumWorkers = 5
)

// job represents a unit of work for the worker pool.
type job struct {
	c *fiber.Ctx // this optional it can bound to other.
}
