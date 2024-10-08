// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gc

import (
	"github.com/valyala/bytebufferpool"
)

// BufferPool is used for efficient memory reuse for I/O operations.
var BufferPool = bytebufferpool.Pool{}
