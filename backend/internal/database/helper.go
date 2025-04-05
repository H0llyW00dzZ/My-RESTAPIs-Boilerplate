// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: The database package here is not covered by tests and won't have tests implemented for it,
// as it is not worth testing the database that requires authentication. (literally stupid testing that requires authentication unlike mock)

package database

import (
	"errors"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// parseDateAdded parses the date_added field from a byte slice.
//
// Note: This helper is useful for MySQL, for example, when you need to convert a time from MySQL into JSON or plain text.
func parseDateAdded(dateAddedBytes []uint8) (time.Time, error) {
	const layout = "2006-01-02 15:04:05" // Define the layout constant
	dateAdded, err := time.Parse(layout, string(dateAddedBytes))
	if err != nil {
		log.LogErrorf("Error parsing date_added: %v", err)
		return time.Time{}, err
	}
	return dateAdded, nil
}

// parseRedisInfo parses the Redis info response and returns a map of key-value pairs.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.SplitSeq(info, "\r\n")
	for line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

// convertStringToInterface converts a slice of strings to a slice of interfaces.
//
// Note: This is pretty useful for big queries, as it can be used with single goroutines or multiple goroutines along with semaphore for MySQL queries.
// Only advanced/master Go developers know how this helper works.
func convertStringToInterface(strs []string) []any {
	// Note: This won't significantly impact performance cost,
	// as it depends on the number of queries (e.g., 1 billion queries will create 1 billion interfaces)
	interfaces := make([]any, len(strs))
	for i, str := range strs {
		interfaces[i] = str
	}
	return interfaces
}

// isDuplicateEntryError checks if an error is a MySQL duplicate entry error.
//
// This function is useful when performing MySQL queries and the goal is to handle duplicate entry errors specifically.
// It takes an error as input and returns a boolean indicating whether the error is a duplicate entry error or not.
//
// Example Usage:
//
//	err := db.Exec("INSERT INTO users (username) VALUES (?)", "gopher")
//	if isDuplicateEntryError(err) {
//	    // Handle duplicate entry error
//	} else if err != nil {
//	    // Handle other errors
//	}
//
// When performing MySQL queries, if an attempt is made to insert a duplicate entry into a unique index or primary key,
// MySQL will return an error with the error number 1062. This function checks if the provided error is a MySQL error
// and if its error number is 1062, indicating a duplicate entry error.
//
// By using this function, duplicate entry errors can be easily identified and handled in MySQL queries without
// the need for string comparisons or manual error code checks.
//
// Note: This function relies on the [github.com/go-sql-driver/mysql] package for the MySQLError type.
// Additionaly Use [WrapMySQLError] function can be used across the codebase to handle various MySQL error types with custom messages.
func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

// IsValidTableName checks if the table name is valid to prevent SQL injection.
func IsValidTableName(name string) bool {
	// Regex to allow alphanumeric characters and underscores, adjust as needed
	validNamePattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(validNamePattern, name)
	if err != nil {
		return false
	}
	return matched
}

// escapeString safely escapes special characters in strings.
//
// This is now correct and can be imported via phpMyAdmin as well.
func escapeString(value string) string {
	value = strings.ReplaceAll(value, "'", "''")
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}

// writeSQLHeader writes a header to the SQL Object.
func writeSQLHeader(w io.Writer) error {
	header := `-- A Better MySQL Dump Written in Go by H0llyW00dzZ
-- Generation Time: ` + time.Now().Format("Jan 02, 2006 at 03:04 PM") + `
` + "\n"
	_, err := w.Write([]byte(header))
	return err
}

// MySQLError represents a custom error type for MySQL errors.
type MySQLError struct {
	OriginalError error
	Number        uint16
	Message       string
}

// Error returns a formatted error message for the MySQLError type.
//
// This method implements the error interface, providing a detailed string
// representation of the MySQL error, including the error number, custom message,
// and the original error.
//
// Example Output:
//
//	"MySQL Error 1062: Error occurred while inserting user - <original error>"
func (e *MySQLError) Error() string {
	return fmt.Sprintf("MySQL Error %d: %s - %v", e.Number, e.Message, e.OriginalError)
}

// WrapMySQLError wraps a MySQL error with a custom message based on the error number.
//
// Example Usage:
//
//	// Example query that might cause a duplicate entry error
//	_, err = db.Exec("INSERT INTO users (username) VALUES (?)", "gopher")
//	if err != nil {
//		err = WrapMySQLError(err, "inserting user")
//		var myErr *MySQLError
//		if errors.As(err, &myErr) {
//			switch myErr.Number {
//			case 1062:
//				// Handle duplicate entry error
//				log.Printf("Duplicate entry detected: %v", myErr)
//			default:
//				// Handle other MySQL errors
//				log.Printf("MySQL error occurred: %v", myErr)
//			}
//		} else {
//			// Handle non-MySQL errors
//			log.Printf("Non-MySQL error occurred: %v", err)
//		}
//	}
//
// Also note that this can be effectively improves HTTP/HTTPS traffic performance as well (e.g., against DDoS for INSERTING Data)
// if MySQL is used directly without Redis (e.g., caching/vice versa).
func WrapMySQLError(err error, context string) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return &MySQLError{
			OriginalError: err,
			Number:        mysqlErr.Number,
			Message:       fmt.Sprintf("Error occurred while %s", context),
		}
	}
	return err
}
