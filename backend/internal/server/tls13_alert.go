// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"crypto/tls"
	"net"
)

// WrappedListener is a custom listener that wraps the TLS listener and handles the RecordHeaderError.
type WrappedListener struct {
	net.Listener
}

// Accept waits for and returns the next connection to the listener.
func (l WrappedListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		if recordHeaderErr, ok := err.(tls.RecordHeaderError); ok {
			return nil, recordHeaderErr
		}
		return nil, err
	}
	return conn, nil
}
