// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package gc provides helper functions and utilities to reduce garbage collector overhead
// and improve memory usage efficiency in Go applications.
//
// The package includes a shared buffer pool for efficient memory reuse during I/O operations,
// string manipulation, and other scenarios involving frequent memory allocations and deallocations.
package gc
