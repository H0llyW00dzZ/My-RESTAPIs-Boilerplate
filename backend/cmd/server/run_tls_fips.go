// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files

// Note: This might needed FIPS-compliant settings (optional).
// Source: https://github.com/microsoft/go/tree/microsoft/main/eng/doc/fips

//go:build tls_fips
// +build tls_fips

package main

import _ "crypto/tls/fipsonly"
