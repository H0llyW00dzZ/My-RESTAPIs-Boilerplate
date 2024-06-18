// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-colorable"
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
func NewLogger(out io.Writer, appName, appNameColor, timeFormat string) *Logger {
	// Enable color output on Windows
	if runtime.GOOS == "windows" {
		out = colorable.NewColorableStdout()
	}

	// Set the desired time format for the loggers
	var flags int
	var timeFormatter func(t time.Time) string

	if timeFormat == TimeFormatUnix {
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

	infoLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%s%s%s] ", appNameColor, appName, ColorReset, ColorGreen, LevelInfo, ColorReset), flags)
	visitorLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%s%s%s] ", appNameColor, appName, ColorReset, ColorYellow, LevelVisitor, ColorReset), flags)
	errorLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%s%s%s] ", appNameColor, appName, ColorReset, ColorRed, LevelError, ColorReset), flags)
	fatalLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%s%s%s] ", appNameColor, appName, ColorReset, ColorMagenta, LevelFatal, ColorReset), flags)
	crashLogger := log.New(NewLogWriter(out, timeFormatter), fmt.Sprintf("[%s%s%s] [%s%s%s] ", appNameColor, appName, ColorReset, ColorBrightRed, LevelCrash, ColorReset), flags)

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
func (l *Logger) Panic(v ...any) {
	l.crashLogger.Print(v...)
}

// Panicf logs formatted panic messages.
func (l *Logger) Panicf(format string, v ...any) {
	l.crashLogger.Printf(format, v...)
}

// Fatal logs fatal messages and exits the program with status code 1.
func (l *Logger) Fatal(v ...any) {
	l.fatalLogger.Fatalln(v...)
}

// Fatalf logs formatted fatal messages and exits the program with status code 1.
func (l *Logger) Fatalf(format string, v ...any) {
	l.fatalLogger.Fatalf(format, v...)
}

// Info logs informational messages.
func (l *Logger) Info(v ...any) {
	l.infoLogger.Println(v...)
}

// Infof logs formatted informational messages.
func (l *Logger) Infof(format string, v ...any) {
	l.infoLogger.Printf(format, v...)
}

// Visitor logs messages related to WebSocket visitors.
func (l *Logger) Visitor(v ...any) {
	l.visitorLogger.Println(v...)
}

// Visitorf logs formatted messages related to WebSocket visitors.
func (l *Logger) Visitorf(format string, v ...any) {
	l.visitorLogger.Printf(format, v...)
}

// Error logs error messages.
func (l *Logger) Error(v ...any) {
	l.errorLogger.Println(v...)
}

// Errorf logs formatted error messages.
func (l *Logger) Errorf(format string, v ...any) {
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
		ColorBlue,  // Color for the application name
		timeFormat, // Time format for the loggers
	)
}

// LogCrash is a convenience function to log crash messages using the custom logger.
func LogCrash(v ...any) {
	customLogger.Panic(v...)
}

// LogCrashf is a convenience function to log formatted crash messages using the custom logger.
func LogCrashf(format string, v ...any) {
	customLogger.Panicf(format, v...)
}

// LogFatal is a convenience function to log fatal messages using the custom logger.
func LogFatal(v ...any) {
	customLogger.Fatal(v...)
}

// LogFatalf is a convenience function to log formatted fatal messages using the custom logger.
func LogFatalf(format string, v ...any) {
	customLogger.Fatalf(format, v...)
}

// LogInfo is a convenience function to log informational messages using the custom logger.
func LogInfo(v ...any) {
	customLogger.Info(v...)
}

// LogInfof is a convenience function to log formatted informational messages using the custom logger.
func LogInfof(format string, v ...any) {
	customLogger.Infof(format, v...)
}

// LogVisitor is a convenience function to log visitor messages using the custom logger.
func LogVisitor(v ...any) {
	customLogger.Visitor(v...)
}

// LogVisitorf is a convenience function to log formatted visitor messages using the custom logger.
func LogVisitorf(format string, v ...any) {
	customLogger.Visitorf(format, v...)
}

// LogError is a convenience function to log error messages using the custom logger.
func LogError(v ...any) {
	customLogger.Error(v...)
}

// LogErrorf is a convenience function to log formatted error messages using the custom logger.
func LogErrorf(format string, v ...any) {
	customLogger.Errorf(format, v...)
}

// LogUserActivity logs a user activity message along with the HTTP method, client IP, User-Agent, and query parameters (if any).
func LogUserActivity(c *fiber.Ctx, activity string) {
	httpMethod := c.Method() // Get the HTTP method of the request
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")
	originalURL := c.OriginalURL() // This gets the full original URL including query parameters
	// TODO: Will implement more features as required. Keep improving, keep coding in Go.

	// Extract the query string from the original URL and format it.
	queryMessage := ""
	if idx := strings.Index(originalURL, "?"); idx != -1 {
		query := originalURL[idx+1:] // Get the substring after the '?' character
		if query != "" {
			queryMessage = fmt.Sprintf(", Query: %s", query)
		}
	}

	// Log the activity with the HTTP method, IP, User-Agent, and query parameters (if any).
	LogVisitorf("Method: %s%s, Activity: %s - IP: %s, User-Agent: %s", httpMethod, queryMessage, activity, clientIP, userAgent)
}
