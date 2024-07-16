// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package logger

// Define ANSI color codes.
const (
	ColorReset     = "\x1b[0m"
	ColorRed       = "\x1b[31m"
	ColorGreen     = "\x1b[32m"
	ColorYellow    = "\x1b[33m"
	ColorBlue      = "\x1b[34m"
	ColorMagenta   = "\x1b[35m"
	ColorBrightRed = "\x1b[91m"
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

// Define Cloudflare formats.
//
// Note: This for internal only, different from frontend.
const (
	CloudflareConnectingIPHeader = "Cf-Connecting-IP"
	UserAgentHeader              = "User-Agent"
	CloudflareRayIDHeader        = "cf-ray"
	CloudflareIPCountryHeader    = "cf-ipcountry"
)
