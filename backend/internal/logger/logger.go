// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger is a custom logger with different levels of logging.
type Logger struct {
	infoLogger    *log.Logger
	visitorLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
	crashLogger   *log.Logger
}

// NewLogger creates a new Logger instance with custom prefixes, colors, and time format.
func NewLogger(out io.Writer, appName, appNameColor, infoColor, visitorColor, errorColor, fatalColor, crashColor, timeFormat string) *Logger {
	// Set the desired time format for the loggers
	var flags int
	var timeFormatter func(t time.Time) string

	if timeFormat == "unix" {
		flags = log.Lmsgprefix
		timeFormatter = func(t time.Time) string {
			return fmt.Sprintf("[%d]", t.Unix())
		}
	} else {
		flags = log.Ldate | log.Ltime | log.Lmsgprefix
		timeFormatter = func(t time.Time) string {
			return ""
		}
	}

	infoLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%sINFO%s] ", appNameColor, appName, ColorReset, infoColor, ColorReset), flags)
	visitorLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%sVISITOR%s] ", appNameColor, appName, ColorReset, visitorColor, ColorReset), flags)
	errorLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%sERROR%s] ", appNameColor, appName, ColorReset, errorColor, ColorReset), flags)
	fatalLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%sFATAL%s] ", appNameColor, appName, ColorReset, fatalColor, ColorReset), flags)
	crashLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%sCRASH%s] ", appNameColor, appName, ColorReset, crashColor, ColorReset), flags)

	return &Logger{
		infoLogger:    infoLogger,
		visitorLogger: visitorLogger,
		errorLogger:   errorLogger,
		fatalLogger:   fatalLogger,
		crashLogger:   crashLogger,
	}
}

// LogWriter is a custom writer that adds the formatted time to the log output.
type LogWriter struct {
	out           io.Writer
	timeFormatter func(t time.Time) string
}

// NewLogWriter creates a new LogWriter instance with the specified output writer and time formatter.
func NewLogWriter(out io.Writer, timeFormatter func(t time.Time) string) *LogWriter {
	return &LogWriter{
		out:           out,
		timeFormatter: timeFormatter,
	}
}

// Write writes the log message with the formatted time to the output writer.
func (w *LogWriter) Write(p []byte) (n int, err error) {
	formattedTime := w.timeFormatter(time.Now())
	if formattedTime != "" {
		formattedTime += " "
	}
	return w.out.Write([]byte(formattedTime + string(p)))
}

// Panic logs panic messages.
func (l *Logger) Panic(v ...interface{}) {
	l.crashLogger.Println(v...)
}

// Panicf logs formatted panic messages.
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.crashLogger.Printf(format, v...)
}

// Fatal logs fatal messages and exits the program with status code 1.
func (l *Logger) Fatal(v ...interface{}) {
	l.fatalLogger.Fatalln(v...)
}

// Fatalf logs formatted fatal messages and exits the program with status code 1.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.fatalLogger.Fatalf(format, v...)
}

// Info logs informational messages.
func (l *Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

// Infof logs formatted informational messages.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Visitor logs messages related to WebSocket visitors.
func (l *Logger) Visitor(v ...interface{}) {
	l.visitorLogger.Println(v...)
}

// Visitorf logs formatted messages related to WebSocket visitors.
func (l *Logger) Visitorf(format string, v ...interface{}) {
	l.visitorLogger.Printf(format, v...)
}

// Error logs error messages.
func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

// Errorf logs formatted error messages.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// customLogger is a package-level variable that will hold the instance of Logger.
// It is initialized using the InitializeLogger function.
var customLogger *Logger

// InitializeLogger sets up the custom logger with the application name and time format.
func InitializeLogger(appName, timeFormat string) {
	customLogger = NewLogger(
		os.Stdout,
		appName,
		ColorBlue,      // Color for the application name
		ColorGreen,     // Color for INFO level
		ColorYellow,    // Color for VISITOR level
		ColorRed,       // Color for ERROR level
		ColorMagenta,   // Color for FATAL level
		ColorBrightRed, // Color for PANIC level
		timeFormat,     // Time format for the loggers
	)
}

// LogCrash is a convenience function to log crash messages using the custom logger.
func LogCrash(v ...interface{}) {
	customLogger.Panic(v...)
}

// LogCrashf is a convenience function to log formatted crash messages using the custom logger.
func LogCrashf(format string, v ...interface{}) {
	customLogger.Panicf(format, v...)
}

// LogFatal is a convenience function to log fatal messages using the custom logger.
func LogFatal(v ...interface{}) {
	customLogger.Fatal(v...)
}

// LogFatalf is a convenience function to log formatted fatal messages using the custom logger.
func LogFatalf(format string, v ...interface{}) {
	customLogger.Fatalf(format, v...)
}

// LogInfo is a convenience function to log informational messages using the custom logger.
func LogInfo(v ...interface{}) {
	customLogger.Info(v...)
}

// LogInfof is a convenience function to log formatted informational messages using the custom logger.
func LogInfof(format string, v ...interface{}) {
	customLogger.Infof(format, v...)
}

// LogVisitor is a convenience function to log visitor messages using the custom logger.
func LogVisitor(v ...interface{}) {
	customLogger.Visitor(v...)
}

// LogVisitorf is a convenience function to log formatted visitor messages using the custom logger.
func LogVisitorf(format string, v ...interface{}) {
	customLogger.Visitorf(format, v...)
}

// LogError is a convenience function to log error messages using the custom logger.
func LogError(v ...interface{}) {
	customLogger.Error(v...)
}

// LogErrorf is a convenience function to log formatted error messages using the custom logger.
func LogErrorf(format string, v ...interface{}) {
	customLogger.Errorf(format, v...)
}

// LogUserActivity logs a user activity message along with the HTTP method, client IP, and User-Agent.
func LogUserActivity(c *fiber.Ctx, activity string) {
	httpMethod := c.Method() // Get the HTTP method of the request
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")

	// Log the activity with the HTTP method, IP, and User-Agent.
	LogVisitorf("Method: %s, Activity: %s - IP: %s, User-Agent: %s", httpMethod, activity, clientIP, userAgent)
}
