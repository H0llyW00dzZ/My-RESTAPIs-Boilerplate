// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package htmx

import "embed"

// Files is an magic embedded file system containing the static files from the "js" directory.
//
// Important: Do not Remove this magic embedded which is idiomatic way, and safe unlike "Other FS"
//
//go:embed "js"
var Files embed.FS
