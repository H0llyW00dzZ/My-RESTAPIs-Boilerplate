// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package worker

import (
	"sync"
	"time"
)

// TokenBucket controls the rate of job processing.
//
// Note: The TokenBucket algorithm is widely used for rate limiting.
// In Kubernetes, it can be enhanced for use with Ingress Nginx while handling
// multiple pods, particularly in scenarios involving Horizontal Pod Autoscaling (HPA).
type TokenBucket struct {
	tokens       int           // Current number of tokens available
	maxTokens    int           // Maximum number of tokens the bucket can hold
	refillRate   time.Duration // Rate at which tokens are added to the bucket
	refillAmount int           // Number of tokens added at each refill interval
	ticker       *time.Ticker  // Ticker to schedule refills
	mu           sync.Mutex    // Mutex to protect concurrent access to the token count
}

// NewTokenBucket initializes a new token bucket with the specified parameters.
//
// This function is designed to be used in conjunction with the worker pool
// to manage concurrency and rate limiting effectively.
//
// TODO: Integrate with [NewDoWork] as this is suitable for managing concurrency.
func NewTokenBucket(maxTokens, refillAmount int, refillRate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		tokens:       maxTokens,
		maxTokens:    maxTokens,
		refillRate:   refillRate,
		refillAmount: refillAmount,
		ticker:       time.NewTicker(refillRate),
	}

	go tb.refill()
	return tb
}

// refill adds tokens to the bucket at the specified rate.
// Ensures that the number of tokens does not exceed the maximum capacity.
func (tb *TokenBucket) refill() {
	for range tb.ticker.C {
		tb.mu.Lock()
		tb.tokens += tb.refillAmount
		if tb.tokens > tb.maxTokens {
			tb.tokens = tb.maxTokens
		}
		tb.mu.Unlock()
	}
}

// Take attempts to take a token from the bucket.
// Returns true if a token was successfully taken, false if no tokens are available.
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// Stop stops the refill ticker, effectively halting the token refill process.
func (tb *TokenBucket) Stop() {
	tb.ticker.Stop()
}
