// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package main

import (
	"h0llyw00dz-template/env"
	"time"
)

// Config holds the application configuration settings
type Config struct {
	AppName         string
	Port            string
	MonitorPath     string
	TimeFormat      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// getConfig retrieves configuration from environment variables or uses default values
func getConfig() Config {
	return Config{
		AppName:         env.GetEnv(env.APPNAME, "Gopher"),
		Port:            env.GetEnv(env.PORT, "8080"),
		MonitorPath:     env.GetEnv(env.MONITORPATH, "/monitor"),
		TimeFormat:      env.GetEnv(env.TIMEFORMAT, "unix"),
		ReadTimeout:     parseDuration(env.GetEnv(env.READTIMEOUT, "5s")),
		WriteTimeout:    parseDuration(env.GetEnv(env.WRITETIMEOUT, "5s")),
		ShutdownTimeout: parseDuration(env.GetEnv(env.SHUTDOWNTIMEOUT, "5s")),
	}
}
