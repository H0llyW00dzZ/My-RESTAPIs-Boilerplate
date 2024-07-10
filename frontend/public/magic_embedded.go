// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package frontend

import "embed"

// Files is an embedded file system containing the static files from the "images, js" directory.
//
// Note: this is a "magic embedded" line and should not be removed, as it is initialized before other code (even another FS System).
//
//go:embed "assets/images"
//go:embed "assets/js"
var Files embed.FS
