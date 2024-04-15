// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package frontend

import "embed"

// Files is an embedded file system containing the static files from the "images" directory.
//
//go:embed "assets/images"
var Files embed.FS
