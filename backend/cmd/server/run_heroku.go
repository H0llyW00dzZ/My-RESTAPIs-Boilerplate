// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: Bubble Tea's spinner must disabled on Heroku because it relies on
// TTY, which is not supported in the Heroku environment. This package is used
// to monitor Go's stack memory usage, garbage collector, and goroutines.

//go:build heroku
// +build heroku

package main

// When this installed, see Docs https://devcenter.heroku.com/articles/language-runtime-metrics-go
// also note that it connected to "run.go", due this method idiom go
import _ "github.com/heroku/x/hmetrics/onload"
