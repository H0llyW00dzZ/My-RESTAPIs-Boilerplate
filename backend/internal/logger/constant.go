// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package logger

// Define ANSI color codes.
const (
	ColorReset     = "\033[0m"
	ColorRed       = "\033[31m"
	ColorGreen     = "\033[32m"
	ColorYellow    = "\033[33m"
	ColorBlue      = "\033[34m"
	ColorMagenta   = "\033[35m"
	ColorBrightRed = "\033[91m"
)

// Define log levels.
const (
	LevelInfo    = "INFO"
	LevelVisitor = "VISITOR"
	LevelError   = "ERROR"
	LevelFatal   = "FATAL"
	LevelCrash   = "CRASH"
)

// Define time formats.
const (
	TimeFormatUnix    = "unix"
	TimeFormatDefault = "default"
)
