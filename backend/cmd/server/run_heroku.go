// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Note: Bubble Tea's spinner must disabled on Heroku because it relies on
// TTY, which is not supported in the Heroku environment. This package is used
// to monitor Go's stack memory usage, garbage collector, and goroutines.

//go:build heroku
// +build heroku

package main

import _ "github.com/heroku/x/hmetrics/onload" // When this installed, see Docs https://devcenter.heroku.com/articles/language-runtime-metrics-go
