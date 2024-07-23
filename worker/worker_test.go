// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker_test

import (
	"context"
	"errors"
	"h0llyw00dz-template/worker"
	"math/rand"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MockJob is a simple job implementation for testing.
type MockJob[T any] struct {
	result    T
	err       error
	sleepTime time.Duration
}

// Execute simulates job execution.
func (j *MockJob[T]) Execute(ctx context.Context) (T, error) {
	if j.sleepTime > 0 {
		time.Sleep(j.sleepTime)
	}
	return j.result, j.err
}

func TestPool_Submit(t *testing.T) {
	pool := worker.NewDoWork[string]()

	// Register a test job
	pool.RegisterJob("testJob", func(c *fiber.Ctx) worker.Job[string] {
		return &MockJob[string]{result: "worker failed to get something from job", err: nil}
	})

	// Submit a job and verify the result
	result, err := pool.Submit(nil, "testJob")
	if err != nil {
		t.Fatalf("Unexpected error during job submission: %v", err)
	}
	if result != "worker failed to get something from job" {
		t.Errorf("Expected result 'worker failed to get something from job', got %s", result)
	}

	// Submit a job that returns an error
	pool.RegisterJob("errorJob", func(c *fiber.Ctx) worker.Job[string] {
		return &MockJob[string]{result: "", err: errors.New("worker failed to get something from job")}
	})

	result, err = pool.Submit(nil, "errorJob")
	if err == nil {
		t.Errorf("Expected error during job submission, but got nil")
	}
	if result != "" {
		t.Errorf("Expected empty result, got %s", result)
	}
	if err.Error() != worker.ErrFailedToGetSomething.Error() {
		t.Errorf("Expected error message 'test error', got %s", err.Error())
	}

	// Submit a job that does not exist
	result, err = pool.Submit(nil, "nonexistentJob")
	if err == nil {
		t.Errorf("Expected error during job submission, but got nil")
	}
	if result != "" {
		t.Errorf("Expected empty result, got %s", result)
	}
	if err.Error() != "job not found: nonexistentJob" {
		t.Errorf("Expected error message 'job not found: nonexistentJob', got %s", err.Error())
	}
}

// Example test for the worker loop
//
// Note: This is a basic example. You should write more comprehensive tests
// covering different scenarios, edge cases, and potential race conditions.
// While this specific implementation uses atomic operations, which are generally
// safe from race conditions, it's still essential to have thorough tests to
// ensure that the code is robust. Additionally, this worker loop has been
// tested in production under high load, handling millions of REST API requests.
//
// Average Memory Consumption (14% of 1GB RAM Pods):
//   - GC Goal Max average: 131MB
//   - Heap average: 120MB
//
// Total:
//   - Latest: 140MB
//   - Average: 164MB
func TestPool_WorkerLoop(t *testing.T) {
	pool := worker.NewDoWork[string]()

	// Register a test job that takes some time to execute
	pool.RegisterJob("slowJob", func(c *fiber.Ctx) worker.Job[string] {
		// Simulate some work (potentially with random delays)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		return &MockJob[string]{result: "slow result", err: nil}
	})

	// Start the pool
	pool.Start()

	// Submit a few jobs
	for i := 0; i < 5; i++ {
		result, err := pool.Submit(nil, "slowJob")
		if err != nil {
			t.Fatalf("Unexpected error during job submission: %v", err)
		}
		if result != "slow result" {
			t.Errorf("Expected result 'slow result', got %s", result)
		}
	}

	// Wait for a short period to allow the jobs to be processed
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Stop the pool
	pool.Stop()

	// Wait for the worker loop to exit (this might be necessary depending on your test setup)
	time.Sleep(worker.DefaultWorkerSleepTime)
}

func TestPool_StartStopLoopZ(t *testing.T) {
	pool := worker.NewDoWork(
		worker.WithNumWorkers[string](1),
		worker.WithJobChannelOptions(worker.WithChanBuffer[worker.Job[string]](1)),
		worker.WithResultChannelOptions(worker.WithChanBuffer[string](1)),
		worker.WithErrorChannelOptions[string](worker.WithChanBuffer[error](1)),
	)

	// Register a test job that takes some time to execute
	pool.RegisterJob("testJob", func(c *fiber.Ctx) worker.Job[string] {
		return &MockJob[string]{result: "test result", sleepTime: time.Millisecond * 50}
	})

	// Submit & Start Job the pool
	pool.Submit(nil, "testJob")

	// Verify that the pool is running
	if !pool.IsRunning() {
		t.Error("Expected pool to be running after job submission")
	}

	// Attempt to start the pool again (should be a no-op)
	pool.Start()

	// Wait for a short period to allow the job to be processed
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	// Stop the pool
	pool.Stop()

	// Wait for the worker loop to exit.
	time.Sleep(worker.DefaultWorkerSleepTime)

	// Verify that the pool is stopped
	if pool.IsRunning() {
		t.Error("Expected pool to be stopped after Stop()")
	}

}
