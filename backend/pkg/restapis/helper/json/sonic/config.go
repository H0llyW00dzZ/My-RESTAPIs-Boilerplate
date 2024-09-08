// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package sonicjson

import (
	"github.com/bytedance/sonic"
)

// ConfigFastest is a pre-configured instance of the [sonic.Config] struct
// with settings optimized for maximum performance.
//
// Note: Use this configuration internally or streaming (which is suitable) for handling JSON data that does not require further re-validation which is already valid json.
// Otherwise, use the default Sonic configuration.
var ConfigFastest = sonic.ConfigFastest
