// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package frontend

import "embed"

// Files is an embedded file system containing the static files from the "images, js, css" directory.
//
// Note: this is a "magic embedded" line and should not be removed, as it is initialized before other code (even another FS System).
//
//go:embed "assets/images"
//go:embed "assets/js"
//go:embed "assets/css"
var Files embed.FS
