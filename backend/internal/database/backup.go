// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package database

import (
	"context"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/gpg"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// BackupTables creates a backup of specified tables in the database.
// It generates a .sql file containing the SQL statements needed to recreate
// the database schema and insert all the current data for each specified table.
//
// TODO: Enhance this function by implementing a background task to automatically back up the database
// using a fully encrypted mechanism provided by the "opengpg/gpg" utility.
// Currently, it only performs a basic backup and stores it directly on disk. Instead,
// the backups should be archived and stored in a cloud storage service.
// This function should be handled by goroutines for task to run periodically without relying on cron jobs,
// as Go supports concurrency with goroutines.
func (s *service) BackupTables(tablesToBackup []string, batchSize int) error {
	for _, tableName := range tablesToBackup {
		if !IsValidTableName(tableName) {
			return fmt.Errorf("invalid table name: %s", tableName)
		}
	}

	// TODO: Implement directory storage for direct disk backup.
	// In the future, this should stream directly with encryption (easy implementation) using OpenPGP/GPG
	// to cloud storage for enhanced security against potential compromises (e.g., between cloud it's self, human error, other).
	backupFile := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.LogErrorf("Failed to close file: %v", cerr)
		}
		if err != nil {
			if remErr := os.Remove(backupFile); remErr != nil {
				log.LogErrorf("Failed to remove incomplete backup file: %v", remErr)
			}
		}
	}()

	// Write the header to the file
	if err = writeSQLHeader(file); err != nil {
		return err
	}

	// For large datasets, this may need to configure this and adjust the MySQL server settings.
	ctx, cancel := context.WithTimeout(context.Background(), DefaultBackupCtxTimeout)
	defer cancel()

	for _, tableName := range tablesToBackup {
		if err = s.dumpTableSchema(ctx, file, tableName); err != nil {
			return err
		}

		if err = s.dumpTableData(ctx, file, tableName, batchSize); err != nil {
			return err
		}
	}

	log.LogInfof("Backup completed: %s", backupFile)
	return nil
}

// BackupTablesConcurrently creates a backup of specified tables concurrently.
// This function leverages goroutines to initiate a backup for each table simultaneously,
// significantly enhancing performance by utilizing multiple CPU cores.
//
// Each table's backup is handled in a separate goroutine, allowing for efficient processing
// of a large number of tables. Errors encountered during the backup process are captured
// via a buffered channel, ensuring that all errors are reported correctly.
//
// Additionally, if this performance improvement is still insufficient for large infrastructures,
// it can be combined with the worker package. Ensure that your infrastructure can handle up to 1 billion operations, as this is just good business.
func (s *service) BackupTablesConcurrently(tablesToBackup []string, o io.Writer, batchSize int) error {
	// Validate table names
	for _, tableName := range tablesToBackup {
		if !IsValidTableName(tableName) {
			return fmt.Errorf("invalid table name: %s", tableName)
		}
	}

	// Write the header once
	// Moved here to prevent it from being written multiple times (Not a Race 🤪). The SQL data remains safe.
	//
	// Note: Ensure there are no errors before processing (on the fly 🛰️). If an error occurs during processing (in the database),
	// the output may be incomplete (e.g., if saving to a file, the header may not be empty).
	// This is due to the real-time nature of the data transfer.
	if err := writeSQLHeader(o); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(tablesToBackup))

	// Start backing up each table concurrently
	for _, tableName := range tablesToBackup {
		wg.Add(1)
		go func(table string) {
			defer wg.Done()
			if err := s.backupSingleTable(table, o, batchSize); err != nil {
				errChan <- err
			}
		}(tableName)
	}

	// Wait for all goroutines to finish and close the error channel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for errors from any of the goroutines
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// BackupTablesWithGPG creates a backup of specified tables in the database and encrypts it using a PGP public key.
func (s *service) BackupTablesWithGPG(tablesToBackup []string, publicKey []string, batchSize int) error {
	for _, tableName := range tablesToBackup {
		if !IsValidTableName(tableName) {
			return fmt.Errorf("invalid table name: %s", tableName)
		}
	}

	// TODO: Implement directory storage for direct disk backup.
	// In the future, this should stream directly with encryption (easy implementation) using OpenPGP/GPG
	// to cloud storage for enhanced security against potential compromises (e.g., between cloud it's self, human error, other).
	backupFile := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	defer func() {
		if file != nil {
			if cerr := file.Close(); cerr != nil {
				log.LogErrorf("Failed to close file: %v", cerr)
			}
		}
		if err != nil {
			if remErr := os.Remove(backupFile); remErr != nil {
				log.LogErrorf("Failed to remove incomplete backup file: %v", remErr)
			}
		}
	}()

	// Write the header to the file
	if err = writeSQLHeader(file); err != nil {
		return err
	}

	// For large datasets, this may need to configure this and adjust the MySQL server settings.
	ctx, cancel := context.WithTimeout(context.Background(), DefaultBackupCtxTimeout)
	defer cancel()

	for _, tableName := range tablesToBackup {
		if err = s.dumpTableSchema(ctx, file, tableName); err != nil {
			return err
		}

		if err = s.dumpTableData(ctx, file, tableName, batchSize); err != nil {
			return err
		}
	}

	// Close the file after writing to it
	if cerr := file.Close(); cerr != nil {
		log.LogErrorf("Failed to close file before encryption: %v", cerr)
		return cerr
	}
	file = nil // Prevent deferred close

	// Encrypt the backup file
	gpg, err := gpg.NewEncryptor(publicKey)
	if err != nil {
		return fmt.Errorf("failed to create encryptor: %w", err)
	}

	encryptedFile := fmt.Sprintf("%s.gpg", backupFile)
	if err = gpg.EncryptFile(backupFile, encryptedFile); err != nil {
		return err
	}

	log.LogInfof("Backup and encryption completed: %s", encryptedFile)

	// Remove the unencrypted backup file
	//
	// Note: It's not possible to remove the file during GPG encryption because EncryptFile streams data using an [io.WriteCloser].
	// This process directly uses a file, not an [io] object, so it should be safe, as well as the file being locked during GPG encryption while open then reading and writing to another file.
	// Consider exploring alternative methods (marked as TODO).
	if err = os.Remove(backupFile); err != nil {
		log.LogErrorf("Failed to remove unencrypted backup file: %v", err)
	}

	return nil
}

// backupSingleTable performs the backup of a single table.
// It creates a pipe to stream the backup data to the specified writer.
// A goroutine is used to write the table's schema and data to the pipe,
// ensuring non-blocking writes. The function validates the table name,
// initiates the backup process, and sends the data to the desired destination.
func (s *service) backupSingleTable(tableName string, o io.Writer, batchSize int) error {
	r, w := io.Pipe()
	defer r.Close()

	go func() {
		defer w.Close()
		if err := s.backupTableToWriter(tableName, w, batchSize); err != nil {
			log.LogErrorf("Failed to backup table %s: %v", tableName, err)
		}
	}()

	// Send the backup data to the desired destination
	if err := s.sendTo(r, o); err != nil {
		return err
	}

	return nil
}

// sendTo transfers data from the reader to the specified writer.
// It uses [io.Copy] to efficiently stream data between the reader and writer,
// allowing for real-time data transfer, such as over a network (on the fly) 🛰️.
//
// TODO: Connect to GPG EncryptStream (laverage I/O before fly)
func (s *service) sendTo(r io.Reader, o io.Writer) error {
	_, err := io.Copy(o, r)
	return err
}

// backupTableToWriter writes the schema and data of a table to the provided writer.
// It begins a transaction to ensure a consistent snapshot of the table.
// The function writes the SQL header, schema, and data to the writer, then commits
// the transaction to finalize the backup.
func (s *service) backupTableToWriter(tableName string, w io.Writer, batchSize int) error {
	// For large datasets, this may need to configure this and adjust the MySQL server settings.
	ctx, cancel := context.WithTimeout(context.Background(), DefaultBackupCtxTimeout)
	defer cancel()

	// Start a transaction to ensure data consistency
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.dumpTableSchema(ctx, w, tableName); err != nil {
		return err
	}

	if err := s.dumpTableData(ctx, w, tableName, batchSize); err != nil {
		return err
	}

	// Commit the transaction after successful backup
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// dumpTableSchema writes the CREATE TABLE statement for the specified table to the object.
func (s *service) dumpTableSchema(ctx context.Context, w io.Writer, tableName string) error {
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	row := s.db.QueryRowContext(ctx, query)
	var table, createTableStmt string
	if err := row.Scan(&table, &createTableStmt); err != nil {
		return fmt.Errorf("failed to get create table statement: %w", err)
	}
	_, err := fmt.Fprintf(w, "%s;\n\n", createTableStmt)
	return err
}

// dumpTableData retrieves all rows from the specified table and writes them as INSERT statements to the object.
//
// Note: This differs from MySQL Dumper and PhpMyAdmin Export, both of which use single-row INSERT statements for data.
// This implementation uses multi-row INSERT statements + Batching, which can improve performance when importing large datasets
// and help avoid MySQL deadlocks (not due to Go, but inherent to MySQL itself).
//
// TODO: Improve batchSize calculations using [math/big].
// The current implementation depends on the architecture (32-bit/64-bit),
// and it should be sufficient for all cases.
func (s *service) dumpTableData(ctx context.Context, w io.Writer, tableName string, batchSize int) error {
	// Adjust the batch size as needed.
	//
	// Note that batching can improve performance when importing data with some MySQL tools,
	// such as MySQL Workbench and PhpMyAdmin. However, in Go, this is often unnecessary
	// because Go can handle streaming efficiently, potentially faster than other languages.
	//
	// This primarily benefits importing, and in Go, it doesn't matter as much due to network factors.
	// For large data (backup), ensure the network is stable, as performance depends on network conditions (e.g., between the client and the MySQL server itself).
	// For example, with large data, you might set the batch size to 100,000 to speed up the import (after backup) process into MySQL.
	//
	// If batchSize is not greater than 0, return an error here.
	// This is a good point unlike return error early (e.g., tableName).
	if batchSize <= 0 {
		return fmt.Errorf("batch size must be greater than 0, got %d", batchSize)
	}

	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query table data: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// This zero allocations because it passes pointers to the values themselves.
	//
	// It supports any Unicode value within the data by using "any" and passing pointers.
	// If there are errors during import after backup, they are typically Unicode issues with external tools.
	// For example, MySQL Workbench might produce errors like "\xF0\x9F\x87\xAE\xF0\x9F..." during import.
	// For large datasets, there might be batch size issues as well.
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var insertStatements []string

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		insertStmt := buildValuesString(values)
		insertStatements = append(insertStatements, insertStmt)

		if len(insertStatements) >= batchSize {
			fullInsert := buildInsertStatement(tableName, columns, insertStatements)
			if _, err := fmt.Fprint(w, fullInsert); err != nil {
				return err
			}
			insertStatements = insertStatements[:0]
		}
	}

	if len(insertStatements) > 0 {
		fullInsert := buildInsertStatement(tableName, columns, insertStatements)
		if _, err := fmt.Fprint(w, fullInsert); err != nil {
			return err
		}
	}
	_, err = fmt.Fprint(w, "\n")
	return err
}

// buildInsertStatement constructs an SQL INSERT statement for multiple row of data.
//
// Note: This differs from MySQL Dumper and PhpMyAdmin Export, both of which use single-row INSERT statements for data.
// This implementation uses multi-row INSERT statements + Batching, which can improve performance when importing large datasets
// and help avoid MySQL deadlocks (not due to Go, but inherent to MySQL itself).
func buildInsertStatement(tableName string, columns []string, values []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO `%s` (", tableName))
	// This is now correct and can be imported via phpMyAdmin as well.
	for i, column := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("`%s`", column))
	}
	sb.WriteString(valuesObject)

	sb.WriteString(strings.Join(values, ", "))
	sb.WriteString(";\n")
	return sb.String()
}

// buildValuesString constructs the VALUES part of an SQL INSERT statement.
func buildValuesString(values []any) string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, val := range values {
		if i > 0 {
			sb.WriteString(", ")
		}
		if val == nil {
			sb.WriteString(nullObject)
		} else if b, ok := val.([]byte); ok {
			sb.WriteString(fmt.Sprintf("'%s'", escapeString(string(b))))
		} else {
			switch v := val.(type) {
			case int64, float64, bool:
				sb.WriteString(fmt.Sprintf("%v", v))
			default:
				sb.WriteString(fmt.Sprintf("'%s'", escapeString(fmt.Sprintf("%v", v))))
			}
		}
	}
	sb.WriteString(")")
	return sb.String()
}
