// Copyright (c) 2025 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package errors_test

import (
	stdErrors "errors"
	"fmt"
	"h0llyw00dz-template/backend/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWrapError tests the Wrap function.
func TestWrapError(t *testing.T) {
	originalErr := stdErrors.New("original error")
	wrappedErr := errors.Wrap(originalErr, "additional context")

	if wrappedErr == nil {
		t.Fatal("expected non-nil error")
	}

	expectedMessage := "additional context: original error"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("expected %q, got %q", expectedMessage, wrappedErr.Error())
	}

	if !stdErrors.Is(wrappedErr, originalErr) {
		t.Error("wrapped error does not contain the original error")
	}
}

// TestWrapErrorF tests the WrapF function.
func TestWrapErrorF(t *testing.T) {
	originalErr := stdErrors.New("original error")
	wrappedErrF := errors.WrapF(originalErr, "additional context with value: %d", 42)

	if wrappedErrF == nil {
		t.Fatal("expected non-nil error")
	}

	expectedMessage := "additional context with value: 42: original error"
	if wrappedErrF.Error() != expectedMessage {
		t.Errorf("expected %q, got %q", expectedMessage, wrappedErrF.Error())
	}

	if !stdErrors.Is(wrappedErrF, originalErr) {
		t.Error("wrapped error does not contain the original error")
	}
}

// TestMultipleWraps tests multiple wrapping of an error.
func TestMultipleWraps(t *testing.T) {
	originalErr := stdErrors.New("original error")
	wrappedErr1 := errors.Wrap(originalErr, "first context")
	wrappedErr2 := errors.Wrap(wrappedErr1, "second context")

	if wrappedErr2 == nil {
		t.Fatal("expected non-nil error")
	}

	expectedMessage := "second context: first context: original error"
	if wrappedErr2.Error() != expectedMessage {
		t.Errorf("expected %q, got %q", expectedMessage, wrappedErr2.Error())
	}

	if !stdErrors.Is(wrappedErr2, originalErr) {
		t.Error("wrapped error does not contain the original error")
	}

	if !stdErrors.Is(wrappedErr2, wrappedErr1) {
		t.Error("wrapped error does not contain the first wrapped error")
	}
}

// TestMultipleWraps10Times tests wrapping an error 10 times.
func TestMultipleWraps10Times(t *testing.T) {
	originalErr := stdErrors.New("original error")
	wrappedErr := originalErr

	// Wrap the error 10 times with different context messages
	for i := 1; i <= 10; i++ {
		wrappedErr = errors.Wrap(wrappedErr, fmt.Sprintf("context %d", i))
	}

	assert.NotNil(t, wrappedErr, "expected non-nil error")

	// Construct the expected message with 10 contexts in chronological order
	expectedMessage := "context 10: context 9: context 8: context 7: context 6: context 5: context 4: context 3: context 2: context 1: original error"

	assert.Equal(t, expectedMessage, wrappedErr.Error(), "unexpected error message")

	// Verify that the wrapped error still contains the original error
	assert.True(t, stdErrors.Is(wrappedErr, originalErr), "wrapped error does not contain the original error")
}

// TestWrapErrorNil tests the Wrap function with a nil error.
func TestWrapErrorNil(t *testing.T) {
	wrappedErr := errors.Wrap(nil, "additional context")
	if wrappedErr != nil {
		t.Fatal("expected nil error")
	}
}

// TestWrapErrorFNil tests the WrapF function with a nil error.
func TestWrapErrorFNil(t *testing.T) {
	wrappedErrF := errors.WrapF(nil, "additional context with value: %d", 42)
	if wrappedErrF != nil {
		t.Fatal("expected nil error")
	}
}

// TestWrapSFNil tests the WrapSF function with a nil error.
func TestWrapSFNil(t *testing.T) {
	wrappedErrSF := errors.WrapF(nil, "additional context with value: %s", "example")
	if wrappedErrSF != nil {
		t.Fatal("expected nil error")
	}
}
