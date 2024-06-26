// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server_test

import (
	"crypto/tls"
	"errors"
	"h0llyw00dz-template/backend/internal/server"
	"net"
	"testing"
)

// mockListener is a mock implementation of net.Listener for testing purposes.
type mockListener struct {
	acceptErr error
}

// Accept returns a mock connection and the configured error.
func (l *mockListener) Accept() (net.Conn, error) {
	return nil, l.acceptErr
}

// Close is a mock implementation of Close method.
func (l *mockListener) Close() error {
	return nil
}

// Addr is a mock implementation of Addr method.
func (l *mockListener) Addr() net.Addr {
	return nil
}

func TestWrappedListener(t *testing.T) {
	t.Run("RecordHeaderError", func(t *testing.T) {
		// Create a mock listener that returns a RecordHeaderError
		recordHeaderErr := tls.RecordHeaderError{Conn: nil, Msg: "record header error"}
		mockListener := &mockListener{acceptErr: recordHeaderErr}

		// Create a wrapped listener with the mock listener
		wrappedListener := server.WrappedListener{Listener: mockListener}

		// Call the Accept method on the wrapped listener
		conn, err := wrappedListener.Accept()

		// Assert that the returned error is the expected RecordHeaderError
		if !errors.As(err, &tls.RecordHeaderError{}) {
			t.Errorf("Expected RecordHeaderError, got: %v", err)
		}

		// Assert that the returned connection is nil
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}

		// Log the error details
		t.Logf("RecordHeaderError: %+v", err)
	})

	t.Run("OtherError", func(t *testing.T) {
		// Create a mock listener that returns a different error
		otherErr := errors.New("other error")
		mockListener := &mockListener{acceptErr: otherErr}

		// Create a wrapped listener with the mock listener
		wrappedListener := server.WrappedListener{Listener: mockListener}

		// Call the Accept method on the wrapped listener
		conn, err := wrappedListener.Accept()

		// Assert that the returned error is the expected error
		if err != otherErr {
			t.Errorf("Expected other error, got: %v", err)
		}

		// Assert that the returned connection is nil
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}

		// Log the error details
		t.Logf("OtherError: %v", err)
	})
}
