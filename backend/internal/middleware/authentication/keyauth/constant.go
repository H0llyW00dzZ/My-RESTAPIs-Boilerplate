// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyauth

import "time"

const (
	sessionKey    = "session"
	apiKey        = "api_key"
	apiKeyExpired = "api_key_expired"
)

const (
	defaultExpryContextKey = time.Second * 2
)
